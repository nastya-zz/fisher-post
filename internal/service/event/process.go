package event

import (
	"context"
	"errors"
	"fmt"

	"post/internal/converter"
	"post/internal/model"
	"post/pkg/logger"
)

var ErrUnknownEventType = errors.New("Process.UnknownEventType")
var ErrNotImplementedEventType = errors.New("Process.NotImplementedEventType")

func (p Processor) Process(ctx context.Context, event model.Event) error {
	const op = "process.Process"

	logger.Info(op, "event", event)

	var user model.User

	switch event.Type {

	case model.UserUpdateProfile:
		user = converter.UserFromPayload(event.Payload)
		p.cachedUserService.Invalidate(ctx, user.ID.String())
		return nil
	case model.UserDelete:
		user = converter.UserFromPayload(event.Payload)
		p.cachedUserService.Invalidate(ctx, user.ID.String())

		err := p.postService.DeleteByUserId(ctx, user.ID)
		if err != nil {
			logger.Error(op, "failed to delete posts", "user id", user.ID, "error", err)
		}
		return nil
	}

	return fmt.Errorf("can't process message %w", ErrUnknownEventType)
}
