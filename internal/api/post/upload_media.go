package post

import (
	"context"
	"fmt"
	"strings"

	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"post/internal/model"
)

func (i *Implementation) UploadMedia(ctx context.Context, req *desc.UploadMediaRequest) (*desc.UploadMediaResponse, error) {
	if err := validateUploadMedia(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	postID, err := model.GetUuid(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id пользователя не валидный")
	}

	createMedia := &model.DescCreateMedia{
		File:     req.Image,
		PostID:   postID,
		Filename: req.Filename,
		Type:     req.Type.String(),
	}

	mediaID, err := i.postService.UploadMedia(ctx, createMedia)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &desc.UploadMediaResponse{
		Id: mediaID.String(),
	}, nil
}

func validateUploadMedia(req *desc.UploadMediaRequest) error {
	var errors []string

	if req.GetPostId() == "" {
		errors = append(errors, "не указан id поста")
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
