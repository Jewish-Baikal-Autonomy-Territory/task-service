package kafka

import (
	"context"
	"errors"
	"fmt"
	domaintask "task-service/internal/domain/task"
	"time"

	analyticspb "task-service/gen/go/task-service/event/analytics"
	taskpb "task-service/gen/go/task-service/task"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreatedEventPublisher struct {
	writer *kafka.Writer
}

func fromDomainCreatedEvent(event *domaintask.CreatedTaskEvent) *analyticspb.TaskCreated {
	task := &analyticspb.TaskCreated{
		TaskId:    new(event.TaskID.String()),
		UserId:    new(event.UserID.String()),
		CreatedAt: timestamppb.New(event.CreatedAt),
	}

	if groupID, ok := event.GroupID.Get(); ok {
		task.GroupId = new(groupID.String())
	}

	if location, ok := event.Location.Get(); ok {
		task.Location = &taskpb.GeoPoint{
			Latitude:  new(location.Latitude),
			Longitude: new(location.Longitude),
		}
	}

	return task
}

func (p *CreatedEventPublisher) Publish(ctx context.Context, event *domaintask.CreatedTaskEvent) error {
	if event == nil {
		return errors.New("event is missing")
	}

	eventPb := fromDomainCreatedEvent(event)
	payload, err := proto.Marshal(eventPb)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(event.UserID.String()),
		Value: payload,
	}
	if err = p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("write messages: %w", err)
	}

	return nil
}

func (p *CreatedEventPublisher) Close() error {
	return p.writer.Close()
}

type TaskCreatedEventPublisherOptions struct {
	Brokers       []string
	Topic         string
	BatchSize     int
	BatchBytes    int64
	BatchTimeout  time.Duration
	RetryAttempts int
}

func (o *TaskCreatedEventPublisherOptions) Validate() error {
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

func NewTaskCreatedEventPublisher(opts TaskCreatedEventPublisherOptions) (*CreatedEventPublisher, error) {
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

	return &CreatedEventPublisher{
		writer: writer,
	}, nil
}
