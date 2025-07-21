package post

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"

	"post/internal/model"
)

// todo: minio client, upload file, create media, return media id
func (s serv) UploadMedia(ctx context.Context, media *model.CreateMedia) (uuid.UUID, error) {
	const op = "post.UploadMedia"
	ext := filepath.Ext(name)

	nameId := uuid.NewString()
	uniqUuidName := fmt.Sprintf("%s-%s%s", media.UserID, nameId, ext)

	link, err := s.minio.UploadFile(ctx, media.File, uniqUuidName)

	// create in repository...

	return uuid.UUID{}, nil
}
