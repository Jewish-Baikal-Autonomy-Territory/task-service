package kafka

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"
	"time"

	eventpb "task-service/gen/go/task-service/event/notification"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type NotificationEventPublisher struct {
	writer *kafka.Writer
}

func fromAppEvent(event *domaintask.NotificationEvent) *eventpb.TaskNotificationInfo {
	return &eventpb.TaskNotificationInfo{
		TaskId: new(event.TaskID.String()),
		UserId: new(event.UserID.String()),
		Title:  new(event.Title),
	}
}

func (p *NotificationEventPublisher) Notify(ctx context.Context, event *domaintask.NotificationEvent) error {
	if event == nil {
		return errors.New("event is missing")
	}

	pbEvent := fromAppEvent(event)
	payload, err := proto.Marshal(pbEvent)

	kafkaMsg := kafka.Message{
		Key:   []byte(event.UserID.String()),
		Value: payload,
	}

	if err = p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return err
	}

	return nil
}

func (p *NotificationEventPublisher) Close() error {
	return p.writer.Close()
}

type NotificationEventPublisherOptions struct {
	Brokers       []string
	Topic         string
	BatchSize     int
	BatchBytes    int64
	BatchTimeout  time.Duration
	RetryAttempts int
}

func (o *NotificationEventPublisherOptions) Validate() error {
	if len(o.Brokers) == 0 {
		return errors.New("brokers is missing")
	}

	if o.Topic == "" {
		return errors.New("topic is missing")
	}

	if o.BatchSize <= 0 {
		return errors.New("batch size is too small")
	}

	if o.BatchBytes <= 0 {
		return errors.New("batch bytes is too small")
	}

	if o.RetryAttempts <= 0 {
		return errors.New("retry attempts is too small")
	}

	return nil
}

func NewTaskNotificationEventPublisher(opts NotificationEventPublisherOptions) (*NotificationEventPublisher, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %w", err)
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(opts.Brokers...),
		Topic:        opts.Topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		BatchSize:    opts.BatchSize,
		BatchBytes:   opts.BatchBytes,
		BatchTimeout: opts.BatchTimeout,
		MaxAttempts:  opts.RetryAttempts,
	}

	return &NotificationEventPublisher{
		writer: writer,
	}, nil
}
