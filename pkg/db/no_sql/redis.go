package no_sql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type Redis struct {
	cli *redis.Client
	ctx context.Context
	logx.Logger
}

type RedisConf struct {
	Addr         string // Redis服务器地址 (host:port)
	Password     string `json:",optional"` // Redis认证密码
	DB           int    `json:",optional"` // 数据库编号
	PoolSize     int    `json:",optional"` // 连接池大小
	MinIdleConns int    `json:",optional"` // 最小空闲连接数
	MaxRetries   int    `json:",optional"` // 最大重试次数
}

// 拼接 Redis 键
func GenerateKey(prefix string, suffix ...string) string {
	parts := []string{prefix}
	parts = append(parts, suffix...)
	return strings.Join(parts, ":")
}

// NewRedisClient 创建一个新的Redis客户端实例
func NewRedisClient(conf RedisConf) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
		MaxRetries:   conf.MaxRetries,
	})
	return &Redis{
		cli:    client,
		ctx:    context.Background(),
		Logger: logx.WithContext(context.Background()),
	}
}

func (r *Redis) Close() error {
	return r.cli.Close()
}

// Set 缓存任意实体对象
func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	if value == nil {
		return fmt.Errorf("value is nil")
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if string(bytes) == "null" || string(bytes) == "[]" {
		return fmt.Errorf("value is nil")
	}

	return r.cli.Set(r.ctx, key, bytes, expiration).Err()
}

// Get 获得缓存的基本对象
func (r *Redis) Get(key string, target interface{}) error {
	bytes, err := r.cli.Get(r.ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

// 锁选项配置
type lockOptions struct {
	LockID    string // 锁持有者唯一标识
	Reentrant bool   // 是否支持可重入
}

// LockOption 定义锁选项函数类型
type LockOption func(*lockOptions)

// WithLockID 指定自定义锁ID（用于追踪锁持有者）
func WithLockID(id string) LockOption {
	return func(o *lockOptions) {
		o.LockID = id
	}
}

// WithReentrant 启用可重入锁模式
func WithReentrant() LockOption {
	return func(o *lockOptions) {
		o.Reentrant = true
	}
}

// LockWithOptions 高级锁方法，支持多种场景配置
// 返回值：(lockID, 是否获取成功, 错误)
func (r *Redis) LockWithOptions(key string, expiration time.Duration, opts ...LockOption) (string, bool, error) {
	options := &lockOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// 生成默认锁ID（如未指定）
	if options.LockID == "" {
		options.LockID = uuid.New().String()
	}

	// 可重入锁逻辑：检查当前锁是否由同一持有者持有
	if options.Reentrant {
		currentVal, err := r.cli.Get(r.ctx, key).Result()
		if err != nil && err != redis.Nil {
			return "", false, fmt.Errorf("检查锁状态失败: %w", err)
		}
		// 已持有锁，延长过期时间并返回成功
		if currentVal == options.LockID {
			if err := r.cli.Expire(r.ctx, key, expiration).Err(); err != nil {
				return options.LockID, false, fmt.Errorf("延长锁过期时间失败: %w", err)
			}
			return options.LockID, true, nil
		}
	}

	// 基础锁逻辑：使用SetNX获取新锁
	acquired, err := r.cli.SetNX(r.ctx, key, options.LockID, expiration).Result()
	if err != nil {
		return options.LockID, false, fmt.Errorf("获取锁失败: %w", err)
	}

	return options.LockID, acquired, nil
}

// 安全解锁Lua脚本（原子操作：检查锁持有者 -> 删除锁）
const unlockScript = `
if redis.call('get', KEYS[1]) == ARGV[1] then
    return redis.call('del', KEYS[1])
else
    return 0
end
`

// UnlockWithID 安全解锁方法（需验证锁持有者ID）
func (r *Redis) UnlockWithID(key, lockID string) error {
	result, err := r.cli.Eval(r.ctx, unlockScript, []string{key}, lockID).Result()
	if err != nil {
		return fmt.Errorf("解锁操作失败: %w", err)
	}
	// 验证脚本执行结果（0表示未解锁成功，1表示成功）
	if res, ok := result.(int64); !ok || res == 0 {
		return fmt.Errorf("解锁失败：锁不存在或持有者不匹配 (lockID=%s)", lockID)
	}
	return nil
}

// Del 删除缓存
func (r *Redis) Del(key string) error {
	return r.cli.Del(r.ctx, key).Err()
}

// DelByPrefix 删除具有指定前缀的所有键
func (r *Redis) DelByPrefix(prefix string) error {
	if prefix == "" {
		return fmt.Errorf("prefix cannot be empty")
	}

	cursor := uint64(0)
	for {
		// 执行 SCAN 命令，返回当前游标和匹配的键
		keys, nextCursor, err := r.cli.Scan(r.ctx, cursor, prefix+"*", 1000).Result()
		if err != nil {
			return fmt.Errorf("SCAN error: %w", err)
		}

		// 使用 goroutines 并发删除键，提高删除速度
		var wg sync.WaitGroup
		for _, key := range keys {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				if err := r.cli.Del(r.ctx, key).Err(); err != nil {
					logx.Errorf("Error deleting key %s: %v", key, err)
				}
			}(key)
		}

		// 等待所有删除操作完成
		wg.Wait()

		// 如果 nextCursor 为 0，说明扫描完成
		if nextCursor == 0 {
			break
		}

		// 更新游标，继续扫描
		cursor = nextCursor
	}

	return nil
}

// Exists 检查缓存是否存在
func (r *Redis) Exists(key string) (bool, error) {
	exists, err := r.cli.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// SetBatch 批量缓存对象
func (r *Redis) SetBatch(data map[string]interface{}, expiration time.Duration) error {
	pipe := r.cli.Pipeline()
	for key, value := range data {
		if value == nil {
			continue
		}
		bytes, err := json.Marshal(value)

		if err != nil {
			return err
		}

		if string(bytes) == "null" || string(bytes) == "[]" {
			continue
		}
		pipe.Set(r.ctx, key, bytes, expiration)
	}
	_, err := pipe.Exec(r.ctx)
	return err
}

// GetBatch 批量获取缓存对象
func (r *Redis) GetBatch(keys []string) (map[string]string, error) {
	result := make(map[string]string)
	pipe := r.cli.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))

	for i, key := range keys {
		cmds[i] = pipe.Get(r.ctx, key)
	}
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return nil, err
	}

	for i, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			continue
		}
		result[keys[i]] = cmd.Val()
	}
	return result, nil
}

// IncrBy 增加缓存值
func (r *Redis) IncrBy(key string, increment int64) (int64, error) {
	return r.cli.IncrBy(r.ctx, key, increment).Result()
}

// DecrBy 减少缓存值
func (r *Redis) DecrBy(key string, decrement int64) (int64, error) {
	return r.cli.DecrBy(r.ctx, key, decrement).Result()
}

// TTL 获取缓存剩余有效时间
func (r *Redis) TTL(key string) (time.Duration, error) {
	return r.cli.TTL(r.ctx, key).Result()
}

// Expire 延长缓存时间
func (r *Redis) Expire(key string, expiration time.Duration) error {
	return r.cli.Expire(r.ctx, key, expiration).Err()
}

// FlushTTLIfLow  如果剩余 TTL 小于 3 小时，则延长缓存时间一天
func (r *Redis) FlushTTLIfLow(key string) error {
	ttl, err := r.TTL(key)
	if err != nil {
		return err
	}

	// 判断剩余时间是否低于 3 小时
	if ttl > 0 && ttl < 3*time.Hour {
		// 延长缓存时间一天
		return r.Expire(key, 24*time.Hour)
	}

	return nil
}

// List操作方法

// LPush 将一个或多个值推入列表的左侧
func (r *Redis) LPush(key string, values ...interface{}) error {
	var strs []string
	for _, value := range values {
		if value == nil {
			continue
		}
		marshal, _ := json.Marshal(value)
		if string(marshal) == "null" || string(marshal) == "[]" {
			continue
		}
		strs = append(strs, string(marshal))
	}

	return r.cli.LPush(r.ctx, key, strs).Err()
}

// LPop 从列表的左侧弹出一个值
func (r *Redis) LPop(key string, model interface{}) error {
	result, err := r.cli.LPop(r.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), model)
}

// RPush 将一个或多个值推入列表的右侧
func (r *Redis) RPush(key string, values ...interface{}) error {
	var strs []string
	for _, value := range values {
		if value == nil {
			continue
		}
		marshal, _ := json.Marshal(value)
		if string(marshal) == "null" || string(marshal) == "[]" {
			continue
		}
		strs = append(strs, string(marshal))
	}

	return r.cli.RPush(r.ctx, key, strs).Err()
}

// RPop 从列表的右侧弹出一个值
func (r *Redis) RPop(key string, model *interface{}) error {
	result, err := r.cli.RPop(r.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), model)
}

// LRange 获取列表中指定范围的元素
func (r *Redis) LRange(key string, start, stop int64) ([]string, error) {
	return r.cli.LRange(r.ctx, key, start, stop).Result()
}

func (r *Redis) LTrim(key string, start, stop int64) (string, error) {
	return r.cli.LTrim(r.ctx, key, start, stop).Result()
}

// Set操作方法

// SAdd 向集合中添加一个或多个成员
func (r *Redis) SAdd(key string, members ...interface{}) error {
	return r.cli.SAdd(r.ctx, key, members...).Err()
}

// SRem 从集合中删除一个或多个成员
func (r *Redis) SRem(key string, members ...interface{}) error {
	return r.cli.SRem(r.ctx, key, members...).Err()
}

// SMembers 获取集合中的所有成员
func (r *Redis) SMembers(key string) ([]string, error) {
	return r.cli.SMembers(r.ctx, key).Result()
}

// SIsMember 判断成员是否在集合中
func (r *Redis) SIsMember(key string, member interface{}) (bool, error) {
	return r.cli.SIsMember(r.ctx, key, member).Result()
}

// GetSetMembers 获取 Redis 集合中的所有成员
func (r *Redis) GetSetMembers(key string) ([]string, error) {
	members, err := r.SMembers(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get members of set %s: %w", key, err)
	}
	return members, nil
}

// Hash操作方法

// HSet 设置哈希表中的字段
func (r *Redis) HSet(key string, values map[string]interface{}) error {
	if values == nil {
		return fmt.Errorf("value is nil")
	}
	return r.cli.HSet(r.ctx, key, values).Err()
}

// HGet 获取哈希表中的字段值
func (r *Redis) HGet(key, field string) (string, error) {
	return r.cli.HGet(r.ctx, key, field).Result()
}

// HGetAll 获取哈希表中的所有字段
func (r *Redis) HGetAll(key string) (map[string]string, error) {
	return r.cli.HGetAll(r.ctx, key).Result()
}

// Z操作方法

// ZAdd 向有序集合中添加一个或多个成员
func (r *Redis) ZAdd(key string, value interface{}, score float64) error {
	if value == nil {
		return fmt.Errorf("value is nil")
	}
	marshal, err := json.Marshal(value)
	if err != nil {
		logx.Errorf("marshal err: %v", err)
		return err
	}

	if string(marshal) == "null" || string(marshal) == "[]" {
		return fmt.Errorf("value is nil")
	}

	z := []redis.Z{
		{
			Score:  score,
			Member: marshal,
		},
	}

	return r.cli.ZAdd(r.ctx, key, z...).Err()
}

// ZRem 从有序集合中删除一个或多个成员
func (r *Redis) ZRem(key string, members ...interface{}) error {
	return r.cli.ZRem(r.ctx, key, members...).Err()
}

// ZRange 获取有序集合中指定范围的成员
func (r *Redis) ZRange(key string, start, stop int64) ([]string, error) {
	return r.cli.ZRange(r.ctx, key, start, stop).Result()
}

// ZRevRange 获取有序集合中指定范围的成员（按分值降序）
func (r *Redis) ZRevRange(key string, start, stop int64) ([]string, error) {
	return r.cli.ZRevRange(r.ctx, key, start, stop).Result()
}

// ZScore 获取指定成员的分值
func (r *Redis) ZScore(key string, member string) (float64, error) {
	return r.cli.ZScore(r.ctx, key, member).Result()
}

// ZRank 获取指定成员的排名（按分值升序）
func (r *Redis) ZRank(key string, member string) (int64, error) {
	return r.cli.ZRank(r.ctx, key, member).Result()
}

// ZRevRank 获取指定成员的排名（按分值降序）
func (r *Redis) ZRevRank(key string, member string) (int64, error) {
	return r.cli.ZRevRank(r.ctx, key, member).Result()
}

// ZCard 获取有序集合的成员数量
func (r *Redis) ZCard(key string) (int64, error) {
	return r.cli.ZCard(r.ctx, key).Result()
}
