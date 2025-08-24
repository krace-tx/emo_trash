package mq

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

// KafkaClient 包含生产者、消费者及管理客户端
type KafkaClient struct {
	brokers       []string             // Kafka broker列表
	config        *sarama.Config       // Kafka配置
	producer      sarama.SyncProducer  // 同步生产者
	asyncProducer sarama.AsyncProducer // 异步生产者
	consumerGroup sarama.ConsumerGroup // 消费者组
	adminClient   sarama.ClusterAdmin  // 集群管理客户端
	mu            sync.Mutex           // 并发安全锁
}

// NewKafkaClient 创建Kafka客户端实例
// brokers: Kafka broker地址列表（如["127.0.0.1:9092"]）
// configOpts: 配置选项函数
func NewKafkaClient(brokers []string, configOpts ...func(*sarama.Config)) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_8_1_0 // 设置Kafka版本（根据实际环境调整）

	// 应用用户自定义配置
	for _, opt := range configOpts {
		opt(config)
	}

	client := &KafkaClient{
		brokers: brokers,
		config:  config,
	}

	// 初始化管理客户端（用于主题管理）
	adminClient, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("初始化管理客户端失败: %w", err)
	}
	client.adminClient = adminClient

	return client, nil
}

// SyncProduce 同步发送消息
// topic: 目标主题
// key: 消息键（用于分区路由）
// value: 消息内容
func (c *KafkaClient) SyncProduce(topic string, key, value []byte) (int32, int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 延迟初始化同步生产者
	if c.producer == nil {
		producer, err := sarama.NewSyncProducer(c.brokers, c.config)
		if err != nil {
			return 0, 0, fmt.Errorf("初始化同步生产者失败: %w", err)
		}
		c.producer = producer
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := c.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, fmt.Errorf("发送消息失败: %w", err)
	}
	return partition, offset, nil
}

// AsyncProduce 异步发送消息
// topic: 目标主题
// key: 消息键
// value: 消息内容
// successChan: 成功回调通道（返回分区和偏移量）
// errorChan: 错误回调通道
func (c *KafkaClient) AsyncProduce(topic string, key, value []byte, successChan chan<- *sarama.ProducerMessage, errorChan chan<- *sarama.ProducerError) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.asyncProducer == nil {
		producer, err := sarama.NewAsyncProducer(c.brokers, c.config)
		if err != nil {
			return fmt.Errorf("初始化异步生产者失败: %w", err)
		}
		c.asyncProducer = producer

		// 启动goroutine处理异步结果
		go func() {
			for {
				select {
				case msg := <-c.asyncProducer.Successes():
					if successChan != nil {
						successChan <- msg
					}
				case err := <-c.asyncProducer.Errors():
					if errorChan != nil {
						errorChan <- err
					}
				}
			}
		}()
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	c.asyncProducer.Input() <- msg
	return nil
}

// Consume 消费消息（消费者组模式）
// ctx: 上下文（用于取消消费）
// groupID: 消费者组ID
// topics: 要消费的主题列表
// handler: 消息处理回调函数
func (c *KafkaClient) Consume(ctx context.Context, groupID string, topics []string, handler func(msg *sarama.ConsumerMessage) error) error {
	c.mu.Lock()
	// 创建消费者组
	consumerGroup, err := sarama.NewConsumerGroup(c.brokers, groupID, c.config)
	if err != nil {
		c.mu.Unlock()
		return fmt.Errorf("初始化消费者组失败: %w", err)
	}
	c.consumerGroup = consumerGroup
	c.mu.Unlock()

	// 实现消费者组处理器
	consumerHandler := &consumerGroupHandler{
		msgHandler: handler,
	}

	// 循环消费（处理再平衡）
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 开始消费
			if err := c.consumerGroup.Consume(ctx, topics, consumerHandler); err != nil {
				return fmt.Errorf("消费失败: %w", err)
			}
		}
	}
}

// CreateTopic 创建Kafka主题
// topic: 主题名称
// partitions: 分区数
// replicationFactor: 副本因子
func (c *KafkaClient) CreateTopic(topic string, partitions int32, replicationFactor int16) error {
	if c.adminClient == nil {
		return fmt.Errorf("管理客户端未初始化")
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	// 创建主题
	err := c.adminClient.CreateTopic(topic, topicDetail, false)
	if err != nil {
		if err == sarama.ErrTopicAlreadyExists {
			return nil // 主题已存在，忽略错误
		}
		return fmt.Errorf("创建主题失败: %w", err)
	}
	return nil
}

// Close 释放Kafka客户端资源
func (c *KafkaClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error

	// 关闭生产者
	if c.producer != nil {
		if e := c.producer.Close(); e != nil {
			err = fmt.Errorf("关闭同步生产者失败: %w", e)
		}
	}

	// 关闭异步生产者
	if c.asyncProducer != nil {
		if e := c.asyncProducer.Close(); e != nil {
			err = fmt.Errorf("%v; 关闭异步生产者失败: %w", err, e)
		}
	}

	// 关闭消费者组
	if c.consumerGroup != nil {
		if e := c.consumerGroup.Close(); e != nil {
			err = fmt.Errorf("%v; 关闭消费者组失败: %w", err, e)
		}
	}

	// 关闭管理客户端
	if c.adminClient != nil {
		if e := c.adminClient.Close(); e != nil {
			err = fmt.Errorf("%v; 关闭管理客户端失败: %w", err, e)
		}
	}

	return err
}

// consumerGroupHandler 实现sarama.ConsumerGroupHandler接口
type consumerGroupHandler struct {
	msgHandler func(msg *sarama.ConsumerMessage) error // 消息处理回调
}

// Setup 消费者组准备就绪时调用
func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup 消费者组关闭时调用
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim 处理分配到的分区消息
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 调用用户提供的消息处理函数
		if err := h.msgHandler(msg); err != nil {
			return err
		}
		// 手动提交偏移量（如果配置为手动提交）
		sess.MarkMessage(msg, "")
	}
	return nil
}
