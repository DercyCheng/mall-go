package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	service "mall-go/services/order-service/application/service"
)

// RabbitMQEventPublisher 基于RabbitMQ的事件发布实现
type RabbitMQEventPublisher struct {
	rabbitMQConn *amqp.Connection
	rabbitMQChan *amqp.Channel
	exchange     string
	logger       *zap.Logger
}

// 确保RabbitMQEventPublisher实现了EventPublisher接口
var _ service.EventPublisher = (*RabbitMQEventPublisher)(nil)

// NewRabbitMQEventPublisher 创建RabbitMQ事件发布器
func NewRabbitMQEventPublisher(
	rabbitMQURL string,
	exchange string,
	logger *zap.Logger,
) (*RabbitMQEventPublisher, error) {
	// 连接到RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	// 创建Channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// 声明Exchange
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQEventPublisher{
		rabbitMQConn: conn,
		rabbitMQChan: ch,
		exchange:     exchange,
		logger:       logger,
	}, nil
}

// EventMessage 事件消息
type EventMessage struct {
	EventName   string      `json:"eventName"`
	EventID     string      `json:"eventId"`
	OccurredOn  time.Time   `json:"occurredOn"`
	Data        interface{} `json:"data"`
	ServiceName string      `json:"serviceName"`
}

// Publish 发布事件
func (p *RabbitMQEventPublisher) Publish(ctx context.Context, eventName string, data interface{}) error {
	// 创建事件消息
	event := EventMessage{
		EventName:   eventName,
		EventID:     generateEventID(),
		OccurredOn:  time.Now(),
		Data:        data,
		ServiceName: "order-service",
	}

	// 序列化事件消息
	body, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event", zap.Error(err), zap.String("eventName", eventName))
		return err
	}

	// 发布消息
	routingKey := "order." + eventName // 使用事件名称作为路由键
	err = p.rabbitMQChan.Publish(
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息
			Timestamp:    time.Now(),
			Headers: amqp.Table{
				"eventName": eventName,
			},
		},
	)

	if err != nil {
		p.logger.Error("Failed to publish event",
			zap.Error(err),
			zap.String("eventName", eventName),
			zap.String("routingKey", routingKey),
		)
		return err
	}

	p.logger.Info("Event published",
		zap.String("eventName", eventName),
		zap.String("routingKey", routingKey),
		zap.String("eventId", event.EventID),
	)
	return nil
}

// Close 关闭连接
func (p *RabbitMQEventPublisher) Close() error {
	if p.rabbitMQChan != nil {
		p.rabbitMQChan.Close()
	}
	if p.rabbitMQConn != nil {
		return p.rabbitMQConn.Close()
	}
	return nil
}

// generateEventID 生成事件ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
		// 增加时间差异，确保随机性
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
