package post

import (
	"context"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"post/internal/model"
	"strconv"
	"strings"
)

func (i *Implementation) CreatePost(ctx context.Context, req *desc.CreatePostRequest) (*desc.Post, error) {

	errors := validation(req)

	if len(errors) > 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s", strings.Join(errors, ", "))
	}

	userID, err := model.GetUuid(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id пользователя не валидный")
	}

	newPost := &model.CreatePost{
		UserID:      userID,
		Description: req.Description,
		Geolocation: model.Geolocation{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		FishTypeIDs:   getIntList(req.FishTypeIds),
		TackleTypeIDs: getIntList(req.TackleTypeIds),
	}

	createdPost, err := i.postService.CreatePost(ctx, newPost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", "Ошибка при создании поста")
	}

	return &desc.Post{
		Id:          createdPost.ID.String(),
		Description: createdPost.Description,
		Location:    req.Location,
		//Media: createdPost.Media, todo convert media
		LikesCount:    int32(createdPost.LikesCount),
		CommentsCount: int32(createdPost.CommentsCount),
		FishTypes:     FromFishTypesToDescFishTypes(createdPost.FishTypes),
		TackleTypes:   FromFishTypesToDescTackleType(createdPost.TackleTypes),
		//CreatedAt: createdPost.CreatedAt, todo convert time

	}, nil
}

func validation(post *desc.CreatePostRequest) []string {
	errors := make([]string, 0, 3)

	media := post.GetMedia()
	if len(media) == 0 {
		errors = append(errors, "Нет медиа файлов")
	}

	userID := post.GetUserId()
	if len(userID) == 0 {
		errors = append(errors, "Не указан id пользователя")
	}

	lat := post.Location.Latitude
	lng := post.Location.Longitude
	if lat < -90.0000000 && lat > 90.0000000 || lng < -180.0000000 && lng > 180.0000000 {
		errors = append(errors, "Координаты заданы вне диапазона")
	}

	return errors
}

func getIntList(list []int32) []int {
	result := make([]int, 0, len(list))

	for _, v := range list {
		result = append(result, int(v))
	}

	return result
}

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
