package mq

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// RocketMQClient RocketMQ客户端结构体
type RocketMQClient struct {
	namesrvAddr  string                // NameServer地址
	producer     rocketmq.Producer     // 生产者实例
	pushConsumer rocketmq.PushConsumer // 推模式消费者
	adminClient  admin.Admin           // 管理客户端
	mu           sync.Mutex            // 并发安全锁
}

// NewRocketMQClient 创建RocketMQ客户端实例
func NewRocketMQClient(namesrvAddr string) *RocketMQClient {
	return &RocketMQClient{
		namesrvAddr: namesrvAddr,
	}
}

// StartProducer 初始化生产者
func (c *RocketMQClient) StartProducer(groupName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.producer != nil {
		return nil // 生产者已初始化
	}

	// 创建生产者实例
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{c.namesrvAddr}),
		producer.WithGroupName(groupName),
		producer.WithRetry(3), // 重试次数
	)
	if err != nil {
		return fmt.Errorf("创建生产者失败: %w", err)
	}

	// 启动生产者
	if err = p.Start(); err != nil {
		return fmt.Errorf("启动生产者失败: %w", err)
	}

	c.producer = p
	return nil
}

// SyncSend 同步发送消息
func (c *RocketMQClient) SyncSend(topic, tag string, body []byte) (*primitive.SendResult, error) {
	if c.producer == nil {
		return nil, fmt.Errorf("生产者未初始化，请先调用StartProducer")
	}

	msg := primitive.NewMessage(topic, body)
	msg.WithTag(tag)

	result, err := c.producer.SendSync(context.Background(), msg)
	if err != nil {
		return nil, fmt.Errorf("同步发送消息失败: %w", err)
	}
	return result, nil
}

// AsyncSend 异步发送消息
func (c *RocketMQClient) AsyncSend(topic, tag string, body []byte, callback func(*primitive.SendResult, error)) error {
	if c.producer == nil {
		return fmt.Errorf("生产者未初始化，请先调用StartProducer")
	}

	msg := primitive.NewMessage(topic, body)
	msg.WithTag(tag)

	err := c.producer.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
		callback(result, err)
	}, msg)

	if err != nil {
		return fmt.Errorf("异步发送消息失败: %w", err)
	}
	return nil
}

// StartPushConsumer 启动推模式消费者
func (c *RocketMQClient) StartPushConsumer(groupName, topic, tag string, handler func(msg *primitive.MessageExt) error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pushConsumer != nil {
		return nil // 消费者已初始化
	}

	// 创建推模式消费者
	pushConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{c.namesrvAddr}),
		consumer.WithGroupName(groupName),
	)
	if err != nil {
		return fmt.Errorf("创建消费者失败: %w", err)
	}

	// 订阅主题
	err = pushConsumer.Subscribe(topic, consumer.MessageSelector{Type: consumer.TAG, Expression: tag},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				if err := handler(msg); err != nil {
					return consumer.ConsumeRetryLater, err
				}
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		return fmt.Errorf("订阅主题失败: %w", err)
	}

	// 启动消费者
	if err = pushConsumer.Start(); err != nil {
		return fmt.Errorf("启动消费者失败: %w", err)
	}

	c.pushConsumer = pushConsumer
	return nil
}

// CreateTopic 创建主题
func (c *RocketMQClient) CreateTopic(topic string, queueNum int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 初始化管理客户端
	if c.adminClient == nil {
		adminClient, err := admin.NewAdmin(admin.WithNamespace(c.namesrvAddr))
		if err != nil {
			return fmt.Errorf("创建管理客户端失败: %w", err)
		}
		c.adminClient = adminClient
	}

	err := c.adminClient.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(topic),
		admin.WithReadQueueNums(queueNum),
		admin.WithWriteQueueNums(queueNum),
	)
	if err != nil {
		return fmt.Errorf("创建主题失败: %w", err)
	}
	return nil
}

// Close 释放资源
func (c *RocketMQClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error

	// 关闭生产者
	if c.producer != nil {
		if e := c.producer.Shutdown(); e != nil {
			err = fmt.Errorf("关闭生产者失败: %w", e)
		}
	}

	// 关闭消费者
	if c.pushConsumer != nil {
		if e := c.pushConsumer.Shutdown(); e != nil {
			if err != nil {
				err = fmt.Errorf("%v; 关闭消费者失败: %w", err, e)
			} else {
				err = fmt.Errorf("关闭消费者失败: %w", e)
			}
		}
	}

	// 关闭管理客户端
	if c.adminClient != nil {
		if e := c.adminClient.Close(); e != nil {
			if err != nil {
				err = fmt.Errorf("%v; 关闭管理客户端失败: %w", err, e)
			} else {
				err = fmt.Errorf("关闭管理客户端失败: %w", e)
			}
		}
	}

	return err
}
