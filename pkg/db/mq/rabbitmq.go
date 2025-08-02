package mq

import (
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"reflect"
	"time"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	url     string
}

// NewRabbitMQ 创建一个新的 RabbitMQ 实例，并支持自动重试连接
func NewRabbitMQ(url, queueName string) *RabbitMQ {
	var conn *amqp.Connection
	var err error
	for retries := 0; retries < 5; retries++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		logx.Errorf("Failed to connect to RabbitMQ (attempt %d): %v", retries+1, err)
		time.Sleep(2 * time.Second) // 重试间隔
	}
	if err != nil {
		logx.Errorf("Failed to connect to RabbitMQ after 5 attempts: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		logx.Errorf("Failed to open a channel: %s", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logx.Errorf("Failed to declare a queue: %s", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		queue:   queue,
		url:     url,
	}
}

// Publish 发送消息到队列
func (r *RabbitMQ) Publish(message any) error {
	var body []byte

	// 判断消息类型并转换为字节数组
	if reflect.TypeOf(message).Kind() != reflect.String {
		body, _ = jsonx.Marshal(message)
	} else {
		body = []byte(message.(string))
	}

	// 使用持久化消息，防止消息丢失
	err := r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 设置消息持久化
		},
	)
	return err
}

// PublishAsync 异步发送消息到队列
func (r *RabbitMQ) PublishAsync(message any, done chan error) {
	go func() {
		err := r.Publish(message)
		done <- err
	}()
}

// Consume 消费队列中的消息，并支持并发处理和错误日志记录
func (r *RabbitMQ) Consume(handler func(string)) {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		false, // auto-ack set to false for manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// 使用协程池或并发处理
	for msg := range msgs {
		go func(msg amqp.Delivery) {
			// 调用处理函数
			handler(string(msg.Body))

			// 手动确认消息
			if err := msg.Ack(false); err != nil {
				logx.Errorf("Failed to acknowledge message: %v", err)
				// 可以记录失败的消息到死信队列等
			}
		}(msg)
	}
}

// GetQueueInfo 获取队列信息
func (r *RabbitMQ) GetQueueInfo() (amqp.Queue, error) {
	queue, err := r.channel.QueueInspect(r.queue.Name)
	if err != nil {
		return amqp.Queue{}, err
	}
	return queue, nil
}

// Close 关闭 RabbitMQ 连接
func (r *RabbitMQ) Close() {
	if err := r.channel.Close(); err != nil {
		log.Fatalf("Failed to close channel: %v", err)
	}
	if err := r.conn.Close(); err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}
