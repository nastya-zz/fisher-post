package converter

import (
	"strconv"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"

	"post/internal/model"
)

func FromFishTypesToDescFishTypes(list []model.Dictionary) []*desc.FishType {
	result := make([]*desc.FishType, 0, len(list))

	for _, v := range list {
		result = append(result, &desc.FishType{
			Id:          strconv.Itoa(v.ID),
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return result
}

func FromFishTypesToDescTackleType(list []model.Dictionary) []*desc.TackleType {
	result := make([]*desc.TackleType, 0, len(list))

	for _, v := range list {
		result = append(result, &desc.TackleType{
			Id:          strconv.Itoa(v.ID),
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return result
}

func FromDescMediaToModelMedia(list []*desc.Media) []*model.Media {
	result := make([]*model.Media, 0, len(list))

	for _, v := range list {
		result = append(result, &model.Media{
			ID:           uuid.MustParse(v.Id),
			MediaType:    v.Type.String(),
			Url:          v.Url,
			ThumbnailUrl: v.ThumbnailUrl,
		})
	}

	return result
}

func FromModelMediaToDescMedia(list []model.Media) []*desc.Media {
	result := make([]*desc.Media, 0, len(list))

	for _, v := range list {
		result = append(result, &desc.Media{
			Id:           v.ID.String(),
			Type:         desc.MediaType(desc.MediaType_value[v.MediaType]),
			Url:          v.Url,
			ThumbnailUrl: v.ThumbnailUrl,
		})
	}

	return result
}

func FromModelUserToDescUser(user model.User) *desc.User {
	return &desc.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		AvatarUrl: user.AvatarUrl,
	}
}
