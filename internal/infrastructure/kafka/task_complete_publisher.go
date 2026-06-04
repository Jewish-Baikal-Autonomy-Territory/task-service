package kafka

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"
	"time"

	analyticspb "task-service/gen/go/task-service/event/analytics"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CompleteTaskPublisher struct {
	writer *kafka.Writer
}

func fromDomainCompleteTaskEvent(event *domaintask.CompleteTaskEvent) *analyticspb.TaskCompleted {
	return &analyticspb.TaskCompleted{
		TaskId:      new(event.TaskID.String()),
		UserId:      new(event.UserID.String()),
		CompletedAt: timestamppb.New(event.CompletedAt),
	}
}

func (p *CompleteTaskPublisher) Publish(ctx context.Context, event *domaintask.CompleteTaskEvent) error {
	if event == nil {
		return errors.New("event is missing")
	}

	pbEvent := fromDomainCompleteTaskEvent(event)
	payload, err := proto.Marshal(pbEvent)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(event.TaskID.String()),
		Value: payload,
	}
	if err = p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("write messages: %w", err)
	}

	return nil
}

func (p *CompleteTaskPublisher) Close() error {
	return p.writer.Close()
}

type CompleteTaskPublisherOptions struct {
	Brokers       []string
	Topic         string
	BatchSize     int
	BatchBytes    int64
	BatchTimeout  time.Duration
	RetryAttempts int
}

func (o *CompleteTaskPublisherOptions) Validate() error {
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

	if o.BatchTimeout <= 0 {
		return errors.New("batch timeout is too small")
	}

	if o.RetryAttempts <= 0 {
		return errors.New("retry attempts is too small")
	}

	return nil
}

func NewCompleteTaskPublisher(opts CompleteTaskPublisherOptions) (*CompleteTaskPublisher, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(opts.Brokers...),
		Topic:        opts.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        true,
		BatchSize:    opts.BatchSize,
		BatchBytes:   opts.BatchBytes,
		BatchTimeout: opts.BatchTimeout,
		MaxAttempts:  opts.RetryAttempts,
	}

	return &CompleteTaskPublisher{
		writer: writer,
	}, nil
}
