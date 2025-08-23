package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// SnowflakeConf 雪花算法配置
type SnowflakeConf struct {
	WorkerID     int64 `json:"worker_id" yaml:"worker_id"`
	DatacenterID int64 `json:"datacenter_id" yaml:"datacenter_id"`
	Epoch        int64 `json:"epoch" yaml:"epoch"`
}

// Snowflake 结构体
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	epoch         int64
	sequence      int64
}

// 定义常量
const (
	workerIDBits     = int64(5)  // 10位中的低5位
	datacenterIDBits = int64(5)  // 10位中的高5位
	sequenceBits     = int64(12) // 序列号位数

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits)     // 最大工作节点ID
	maxDatacenterID = int64(-1) ^ (int64(-1) << datacenterIDBits) // 最大数据中心ID
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)     // 最大序列号

	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
)

// NewSnowflakeFromConfig 从配置文件创建Snowflake实例
func NewSnowflakeFromConfig(config *SnowflakeConf) (*Snowflake, error) {

	if config.WorkerID < 0 || config.WorkerID > maxWorkerID {
		return nil, errors.New(fmt.Sprintf("worker ID must be between 0 and %d", maxWorkerID))
	}
	if config.DatacenterID < 0 || config.DatacenterID > maxDatacenterID {
		return nil, errors.New(fmt.Sprintf("datacenter ID must be between 0 and %d", maxDatacenterID))
	}

	// 如果没有配置epoch，使用默认值
	epoch := config.Epoch
	if epoch == 0 {
		epoch = 1609459200000 // 2021-01-01 00:00:00 UTC
	}

	return &Snowflake{
		lastTimestamp: 0,
		workerID:      config.WorkerID,
		datacenterID:  config.DatacenterID,
		epoch:         epoch,
		sequence:      0,
	}, nil
}

// MustNewSnowflake 从配置文件创建Snowflake实例，若失败则panic
func MustNewSnowflake(config SnowflakeConf) *Snowflake {
	snowflake, err := NewSnowflake(config)
	if err != nil {
		panic(err)
	}
	return snowflake
}

// NewSnowflake 直接创建Snowflake实例
func NewSnowflake(config SnowflakeConf) (*Snowflake, error) {
	if config.WorkerID < 0 || config.WorkerID > maxWorkerID {
		return nil, errors.New(fmt.Sprintf("worker ID must be between 0 and %d", maxWorkerID))
	}
	if config.DatacenterID < 0 || config.DatacenterID > maxDatacenterID {
		return nil, errors.New(fmt.Sprintf("datacenter ID must be between 0 and %d", maxDatacenterID))
	}

	if config.Epoch == 0 {
		config.Epoch = 1609459200000
	}

	return &Snowflake{
		lastTimestamp: 0,
		workerID:      config.WorkerID,
		datacenterID:  config.DatacenterID,
		epoch:         config.Epoch,
		sequence:      0,
	}, nil
}

// Generate 生成一个唯一的ID
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1000000 // 转换为毫秒

	if timestamp < s.lastTimestamp {
		panic("clock moved backwards")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 当前毫秒内的序列号已用完，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	// 组合各部分生成ID
	id := ((timestamp - s.epoch) << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

// ParseID 解析雪花ID
func (s *Snowflake) ParseID(id int64) (timestamp int64, workerID int64, datacenterID int64, sequence int64) {
	timestamp = (id >> timestampLeftShift) + s.epoch
	datacenterID = (id >> datacenterIDShift) & maxDatacenterID
	workerID = (id >> workerIDShift) & maxWorkerID
	sequence = id & maxSequence
	return
}

// GetConfigInfo 获取配置信息
func (s *Snowflake) GetConfigInfo() (workerID, datacenterID, epoch int64) {
	return s.workerID, s.datacenterID, s.epoch
}
