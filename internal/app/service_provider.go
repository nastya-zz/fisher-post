package app

import (
	"context"

	"post/internal/api/post"
	"post/internal/client/db"
	"post/internal/client/db/pg"
	userservice "post/internal/client/user_service"
	"post/internal/closer"
	"post/internal/config"
	"post/internal/repository"
	commentRepository "post/internal/repository/comment"
	likeRepository "post/internal/repository/like"
	postRepository "post/internal/repository/post"
	"post/internal/service"
	commentService "post/internal/service/comment"
	likeService "post/internal/service/like"
	postService "post/internal/service/post"
	"post/internal/transaction"
	"post/pkg/logger"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	rmqConfig  config.RMQConfig

	//rmqClient       broker.ClientMsgBroker
	dbClient          db.Client
	userServiceClient userservice.ServiceClient
	txManager         db.TxManager
	postRepository    repository.PostRepository
	commentRepository repository.CommentRepository
	likeRepository    repository.LikeRepository

	postService    service.PostService
	commentService service.CommentService
	likeService    service.LikeService
	postImpl       *post.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			logger.Fatal("failed to get pg config", "error", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) RMQConfig() config.RMQConfig {
	if s.rmqConfig == nil {
		cfg, err := config.NewRMQConfig()
		if err != nil {
			logger.Fatal("failed to get rmqConfig", "error", err.Error())
		}

		s.rmqConfig = cfg
	}

	return s.rmqConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			logger.Fatal("failed to get grpc config", "error", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Fatal("failed to create db client", "error", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			logger.Fatal("ping error", "error", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UserServiceClient(ctx context.Context) userservice.ServiceClient {
	if s.userServiceClient == nil {
		userServiceClient, err := userservice.New(ctx)
		if err != nil {
			logger.Fatal("failed to create user service client", "error", err)
		}
		closer.Add(userServiceClient.Close)
		s.userServiceClient = userServiceClient
	}

	return s.userServiceClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) PostRepository(ctx context.Context) repository.PostRepository {
	if s.postRepository == nil {
		s.postRepository = postRepository.New(s.DBClient(ctx))
	}

	return s.postRepository
}

func (s *serviceProvider) LikeRepository(ctx context.Context) repository.LikeRepository {
	if s.likeRepository == nil {
		s.likeRepository = likeRepository.New(s.DBClient(ctx))
	}

	return s.likeRepository
}

func (s *serviceProvider) CommentRepository(ctx context.Context) repository.CommentRepository {
	if s.commentRepository == nil {
		s.commentRepository = commentRepository.New(s.DBClient(ctx))
	}

	return s.commentRepository
}

func (s *serviceProvider) PostService(ctx context.Context) service.PostService {
	if s.postService == nil {
		s.postService = postService.New(
			s.PostRepository(ctx),
			s.TxManager(ctx),
			s.UserServiceClient(ctx),
			s.LikeRepository(ctx),
		)
	}

	return s.postService
}

func (s *serviceProvider) CommentService(ctx context.Context) service.CommentService {
	if s.commentService == nil {
		s.commentService = commentService.New(
			s.CommentRepository(ctx),
		)
	}

	return s.commentService
}

func (s *serviceProvider) LikeService(ctx context.Context) service.LikeService {

	if s.likeService == nil {
		s.likeService = likeService.New(s.LikeRepository(ctx))
	}

	return s.likeService
}

func (s *serviceProvider) PostImpl(ctx context.Context) *post.Implementation {

	if s.postImpl == nil {
		s.postImpl = post.NewImplementation(s.PostService(ctx), s.CommentService(ctx), s.LikeService(ctx))
	}

	return s.postImpl
}
