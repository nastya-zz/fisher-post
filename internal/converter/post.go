package converter

import (
	"strconv"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/protobuf/types/known/timestamppb"

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

func FromModelPostToDescPost(post *model.Post) *desc.Post {
	return &desc.Post{
		Id:          post.ID.String(),
		User:        FromModelUserToDescUser(post.User),
		Description: post.Description,
		Location: &desc.LatLng{
			Latitude:  post.Geolocation.Latitude,
			Longitude: post.Geolocation.Longitude,
		},
		Media:         FromModelMediaToDescMedia(post.Media),
		LikesCount:    int32(post.LikesCount),
		CommentsCount: int32(post.CommentsCount),
		FishTypes:     FromFishTypesToDescFishTypes(post.FishTypes),
		TackleTypes:   FromFishTypesToDescTackleType(post.TackleTypes),
		CreatedAt:     timestamppb.New(post.CreatedAt),
	}
}

func FromModelCommentToDescComment(comment *model.Comment) *desc.Comment {
	// Конвертация ParentCommentID в *string
	var parentCommentId *string
	if comment.ParentCommentID != nil {
		parentCommentIdStr := comment.ParentCommentID.String()
		parentCommentId = &parentCommentIdStr
	}

	// Конвертация ReplyToUserID в *string
	var replyToUserId *string
	if comment.ReplyToUserID != nil {
		replyToUserIdStr := comment.ReplyToUserID.String()
		replyToUserId = &replyToUserIdStr
	}

	// Рекурсивная конвертация ответов
	var replies []*desc.Comment
	if comment.Replies != nil {
		replies = FromModelCommentsToDescComments(comment.Replies)
	}

	return &desc.Comment{
		Id:              comment.ID.String(),
		User:            FromModelUserToDescUser(*comment.User),
		PostId:          comment.PostID.String(),
		Content:         comment.Content,
		CreatedAt:       timestamppb.New(comment.CreatedAt),
		UpdatedAt:       timestamppb.New(comment.UpdatedAt),
		UserId:          comment.UserID.String(),
		ParentCommentId: parentCommentId,
		ReplyToUserId:   replyToUserId,
		Replies:         replies,
		IsReply:         comment.ParentCommentID != nil,
	}
}

func FromModelCommentsToDescComments(comments []*model.Comment) []*desc.Comment {
	result := make([]*desc.Comment, 0, len(comments))
	for _, comment := range comments {
		result = append(result, FromModelCommentToDescComment(comment))
	}
	return result
}
