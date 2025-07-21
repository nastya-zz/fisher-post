package post

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"post/internal/model"
)

func (i *Implementation) UploadMedia(ctx context.Context, req *desc.UploadMediaRequest) (*desc.UploadMediaResponse, error) {
	if err := validateUploadMedia(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	userID, err := model.GetUuid(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id пользователя не валидный")
	}

	i.service.UploadMedia(ctx, userID, req.Filename, req.Type)

	return &desc.UploadMediaResponse{
		Id:     uuid.New().String(),
	}, nil
}

func validateUploadMedia(req *desc.UploadMediaRequest) error {
	var errors []string

	if req.GetUserId() == "" {
		errors = append(errors, "не указан id пользователя")
	}

	if req.GetFilename() == "" {
		errors = append(errors, "файл не загружен")
	}

	mediaType := strings.ToUpper(req.GetType().String())
	if mediaType != model.MediaTypePhoto && mediaType != model.MediaTypeVideo {
		errors = append(errors, "неверный тип медиа")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}
