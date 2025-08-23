package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

var (
	RabbitMQConn *amqp.Connection
	RabbitMQChan *amqp.Channel
)

// UserBehaviorEvent 用户行为事件
type UserBehaviorEvent struct {
	UserID       int64                  `json:"user_id"`
	VideoID      int64                  `json:"video_id"`
	BehaviorType string                 `json:"behavior_type"` // watch, like, comment, share
	Duration     int32                  `json:"duration"`      // 观看时长
	Timestamp    time.Time              `json:"timestamp"`
	ExtraData    map[string]interface{} `json:"extra_data,omitempty"`
}

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ(url string) error {
	var err error
	RabbitMQConn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	RabbitMQChan, err = RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}

	// 声明交换机
	err = RabbitMQChan.ExchangeDeclare(
		"user_behavior", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// 声明队列
	_, err = RabbitMQChan.QueueDeclare(
		"user_behavior_queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// 绑定队列到交换机
	err = RabbitMQChan.QueueBind(
		"user_behavior_queue", // queue name
		"user.behavior.*",     // routing key
		"user_behavior",       // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	klog.Info("RabbitMQ initialized successfully")
	return nil
}

// PublishUserBehavior 发布用户行为事件
func PublishUserBehavior(ctx context.Context, event *UserBehaviorEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	routingKey := fmt.Sprintf("user.behavior.%s", event.BehaviorType)

	err = RabbitMQChan.PublishWithContext(ctx,
		"user_behavior", // exchange
		routingKey,      // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 消息持久化
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	klog.Infof("Published user behavior event: user_id=%d, video_id=%d, type=%s",
		event.UserID, event.VideoID, event.BehaviorType)
	return nil
}

// ConsumeUserBehavior 消费用户行为事件
func ConsumeUserBehavior(handler func(*UserBehaviorEvent) error) error {
	msgs, err := RabbitMQChan.Consume(
		"user_behavior_queue", // queue
		"",                    // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var event UserBehaviorEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				klog.Errorf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := handler(&event); err != nil {
				klog.Errorf("Failed to handle event: %v", err)
				msg.Nack(false, true) // 重新入队
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

// Close 关闭连接
func Close() {
	if RabbitMQChan != nil {
		RabbitMQChan.Close()
	}
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
}
