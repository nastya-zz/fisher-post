package post

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"

	"post/internal/model"
)

func (s serv) UploadMedia(ctx context.Context, media *model.DescCreateMedia) (uuid.UUID, error) {
	const op = "post.UploadMedia"
	ext := filepath.Ext(media.Filename)

	nameId := uuid.NewString()
	uniqUuidName := fmt.Sprintf("%s-%s%s", media.PostID, nameId, ext)

	link, err := s.minioService.UploadFile(ctx, media.File, uniqUuidName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	mediaForSave := &model.CreateMedia{
		PostID: media.PostID,
		URL: link,
		Type: media.Type,
	}

	if media.IsThumbnail {
		mediaForSave.ThumbnailURL = link
	}

	mediaID, err := s.mediaRepository.Create(ctx, mediaForSave)	
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return mediaID, nil
}
