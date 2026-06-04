package kafka

import (
	"context"
	"errors"
	"fmt"
	pbsystem "task-service/gen/go/task-service/event/system"
	appevent "task-service/internal/application/event"
	"time"

	"buf.build/go/protovalidate"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type UserDeletedConsumer struct {
	reader    *kafka.Reader
	handler   appevent.DeleteUserTasksHandler
	validator protovalidate.Validator
	logger    *zerolog.Logger
}

func (c *UserDeletedConsumer) Handle(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("fetch message: %w", err)
		}

		if err != nil {
			c.logger.Error().
				Str("UserDeletedConsumer", "Handle").
				Err(err).
				Msg("error reading message")
			continue
		}

		var userDeleted pbsystem.UserDeleted
		if err = proto.Unmarshal(msg.Value, &userDeleted); err != nil {
			c.logger.Error().
				Str("UserDeletedConsumer", "Handle").
				Err(err).
				Msg("error unmarshalling message")
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		if err = c.validator.Validate(&userDeleted); err != nil {
			c.logger.Error().
				Str("UserDeletedConsumer", "Handle").
				Err(err).
				Msg("error validating user deleted event")
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		userID, err := uuid.Parse(userDeleted.GetUserId())
		if err != nil {
			c.logger.Error().
				Str("UserDeletedConsumer", "Handle").
				Err(err).
				Msg("error parsing user id")
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		if err = c.handler.Handle(ctx, userID); err != nil {
			c.logger.Error().
				Str("UserDeletedConsumer", "Handle").
				Err(err).
				Msg("error processing user deleted event")
			continue
		}

		if err = c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error().
				Str("UserDeletedConsumer", "Handle").
				Err(err).
				Msg("error committing user deleted event")
		}
	}
}

func (c *UserDeletedConsumer) Stop() error {
	return c.reader.Close()
}

type UserDeletedConsumerOptions struct {
	Logger            *zerolog.Logger
	Brokers           []string
	GroupID           string
	Topic             string
	MinBytes          int
	MaxBytes          int
	MaxWait           time.Duration
	HeartbeatInterval time.Duration
	SessionTimeout    time.Duration
}

func (o *UserDeletedConsumerOptions) Validate() error {
	if len(o.Brokers) == 0 {
		return errors.New("missing brokers")
	}

	if o.GroupID == "" {
		return errors.New("missing group id")
	}

	if o.Topic == "" {
		return errors.New("missing topic")
	}

	if o.MinBytes == 0 {
		return errors.New("missing min bytes")
	}

	if o.MaxBytes < o.MinBytes {
		return errors.New("max bytes too small")
	}

	if o.MaxWait == 0 {
		return errors.New("max wait too small")
	}

	if o.HeartbeatInterval == 0 {
		return errors.New("heartbeat interval too small")
	}

	if o.SessionTimeout == 0 {
		return errors.New("session timeout too small")
	}

	return nil
}

func NewUserDeletedConsumer(h appevent.DeleteUserTasksHandler, opts UserDeletedConsumerOptions) (*UserDeletedConsumer, error) {
	if h == nil {
		return nil, errors.New("handler is missing")
	}

	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %w", err)
	}

	validator, err := protovalidate.New(protovalidate.WithFailFast())
	if err != nil {
		return nil, errors.New("validator is invalid")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           opts.Brokers,
		GroupID:           opts.GroupID,
		Topic:             opts.Topic,
		MinBytes:          opts.MinBytes,
		MaxBytes:          opts.MaxBytes,
		MaxWait:           opts.MaxWait,
		HeartbeatInterval: opts.HeartbeatInterval,
		SessionTimeout:    opts.SessionTimeout,
		StartOffset:       kafka.FirstOffset,
	})

	return &UserDeletedConsumer{
		reader:    reader,
		handler:   h,
		validator: validator,
		logger:    opts.Logger,
	}, nil
}
