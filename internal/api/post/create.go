package post

import (
	"context"
	"strings"

	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"post/internal/converter"
	"post/internal/model"
	"post/pkg/logger"
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
		logger.Error("failed to create post", "error", err)
		return nil, status.Errorf(codes.Internal, "%s", "Ошибка при создании поста")
	}

	return &desc.Post{
		Id:            createdPost.ID.String(),
		User:          converter.FromModelUserToDescUser(createdPost.User),
		Description:   createdPost.Description,
		Location:      req.Location,
		Media:         converter.FromModelMediaToDescMedia(createdPost.Media),
		LikesCount:    int32(createdPost.LikesCount),
		CommentsCount: int32(createdPost.CommentsCount),
		FishTypes:     converter.FromFishTypesToDescFishTypes(createdPost.FishTypes),
		TackleTypes:   converter.FromFishTypesToDescTackleType(createdPost.TackleTypes),
		CreatedAt:     timestamppb.New(createdPost.CreatedAt),
	}, nil
}

func validation(post *desc.CreatePostRequest) []string {
	errors := make([]string, 0, 3)

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
