package app

import (
	"context"

	"post/internal/api/post"
	broker "post/internal/client/broker"
	"post/internal/client/broker/rabbitmq"
	cache "post/internal/client/cache"
	redis "post/internal/client/cache/redis"
	"post/internal/client/cache/redis/invalidators"
	"post/internal/client/db"
	"post/internal/client/db/pg"
	"post/internal/client/minio/minio"
	userservice "post/internal/client/user_service"
	"post/internal/closer"
	"post/internal/config"
	"post/internal/consumer"
	rmqConsumer "post/internal/consumer/rabbitmq"
	rmqProducer "post/internal/producer/rabbitmq"
	"post/internal/repository"
	commentRepository "post/internal/repository/comment"
	eventRepository "post/internal/repository/event"
	likeRepository "post/internal/repository/like"
	mediaRepository "post/internal/repository/media"
	postRepository "post/internal/repository/post"
	"post/internal/service"
	casheduserservice "post/internal/service/cashed_user_service"
	commentService "post/internal/service/comment"
	"post/internal/service/event"
	likeService "post/internal/service/like"
	minioService "post/internal/service/minio"
	postService "post/internal/service/post"
	userService "post/internal/service/user"
	"post/internal/transaction"
	"post/pkg/logger"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	rmqConfig     config.RMQConfig
	minioConfig   *config.MinioConfig
	redisConfig   config.RedisConfig
	eventConsumer consumer.Consumer

	postProducer rmqProducer.Producer

	rmqClient         broker.ClientMsgBroker
	dbClient          db.Client
	userServiceClient userservice.ServiceClient
	minioClient       *minio.Client
	txManager         db.TxManager
	postRepository    repository.PostRepository
	commentRepository repository.CommentRepository
	likeRepository    repository.LikeRepository
	mediaRepository   repository.MediaRepository
	eventRepository   repository.EventRepository

	postService       service.PostService
	commentService    service.CommentService
	likeService       service.LikeService
	postImpl          *post.Implementation
	minioService      service.MinioService
	cachedUserService service.CachedUserService
	cacheService      cache.Cache
	userService       service.UserService
	eventService      service.EventsService
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

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			logger.Fatal("failed to get redis config", "error", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
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

func (s *serviceProvider) MinioConfig() config.MinioConfig {
	if s.minioConfig == nil {
		cfg, err := config.NewMinioConfig()
		if err != nil {
			logger.Fatal("failed to get minio config", "error", err.Error())
		}

		s.minioConfig = cfg
	}

	return *s.minioConfig
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

func (s *serviceProvider) MinioClient(ctx context.Context) *minio.Client {
	if s.minioClient == nil {
		cfg := s.MinioConfig()
		cl, err := minio.New(ctx, cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
		if err != nil {
			logger.Fatal("failed to create minio client", "error", err)
		}

		s.minioClient = cl
	}

	return s.minioClient
}

func (s *serviceProvider) RabbitMQClient(ctx context.Context) broker.ClientMsgBroker {
	if s.rmqClient == nil {
		cl, err := rabbitmq.NewRabbitMQ(ctx, s.RMQConfig().DSN())
		if err != nil {
			logger.Fatal("failed to create rmq client", "error", err)
		}

		closer.Add(cl.Close)

		s.rmqClient = cl
	}

	return s.rmqClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) EventRepository(ctx context.Context) repository.EventRepository {
	if s.eventRepository == nil {
		s.eventRepository = eventRepository.NewRepository(s.DBClient(ctx))
	}

	return s.eventRepository
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

func (s *serviceProvider) MediaRepository(ctx context.Context) repository.MediaRepository {

	if s.mediaRepository == nil {
		s.mediaRepository = mediaRepository.New(s.DBClient(ctx))
	}

	return s.mediaRepository
}

func (s *serviceProvider) PostService(ctx context.Context) service.PostService {
	if s.postService == nil {
		s.postService = postService.New(
			s.PostRepository(ctx),
			s.TxManager(ctx),
			s.GetCachedUserService(ctx),
			s.LikeRepository(ctx),
			s.MediaRepository(ctx),
			s.MinioService(ctx),
			s.EventRepository(ctx),
		)
	}

	return s.postService
}
func (s *serviceProvider) UserService(ctx context.Context, userClient userservice.ServiceClient) service.UserService {
	if s.userService == nil {
		s.userService = userService.New(
			s.UserServiceClient(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) CommentService(ctx context.Context) service.CommentService {
	if s.commentService == nil {
		s.commentService = commentService.New(
			s.CommentRepository(ctx),
			s.GetCachedUserService(ctx),
		)
	}

	return s.commentService
}

func (s *serviceProvider) LikeService(ctx context.Context) service.LikeService {

	if s.likeService == nil {
		s.likeService = likeService.New(s.LikeRepository(ctx), s.GetCachedUserService(ctx))
	}

	return s.likeService
}

func (s *serviceProvider) MinioService(ctx context.Context) service.MinioService {

	if s.minioService == nil {
		s.minioService = minioService.New(s.MinioClient(ctx))
	}

	return s.minioService
}

func (s *serviceProvider) PostImpl(ctx context.Context) *post.Implementation {

	if s.postImpl == nil {
		s.postImpl = post.NewImplementation(s.PostService(ctx), s.CommentService(ctx), s.LikeService(ctx))
	}

	return s.postImpl
}

func (s *serviceProvider) GetCachedUserService(ctx context.Context) service.CachedUserService {
	if s.cachedUserService == nil {
		// Базовый сервис
		baseUserService := s.UserService(ctx, s.UserServiceClient(ctx))

		// Кеш
		cache := s.GetCacheService(ctx)

		// Невалидатор
		invalidator := invalidators.NewCacheInvalidator(cache)

		// Декоратор с кешированием
		s.cachedUserService = casheduserservice.New(baseUserService, cache, invalidator)
	}

	return s.cachedUserService
}

func (s *serviceProvider) GetCacheService(ctx context.Context) cache.Cache {
	if s.cacheService == nil {
		cfg := s.RedisConfig()
		s.cacheService = redis.NewRedisCache(
			cfg.Config().Addr,
			cfg.Config().Password,
			cfg.Config().DB,
		)
		err := s.cacheService.Ping(ctx)
		if err != nil {
			logger.Fatal("failed to ping redis", "error", err.Error())
		}
		closer.Add(s.cacheService.Close)
	}

	return s.cacheService
}

func (s *serviceProvider) PostProducer(ctx context.Context) rmqProducer.Producer {
	if s.postProducer == nil {
		postProducer, err := rmqProducer.NewPostProducer(s.RabbitMQClient(ctx).Connect())
		if err != nil {
			logger.Fatal("failed to create post producer", "error", err.Error())
		}
		// Ensure type compatibility; store the producer as the interface
		s.postProducer = postProducer
	}
	return s.postProducer
}

func (s *serviceProvider) EventService(ctx context.Context) service.EventsService {
	if s.eventService == nil {
		s.eventService = event.New(
			s.GetCachedUserService(ctx),
			s.PostService(ctx),
			s.PostProducer(ctx),
			s.EventRepository(ctx),
		)
	}

	return s.eventService
}

func (s *serviceProvider) EventConsumer(ctx context.Context) consumer.Consumer {
	if s.eventConsumer == nil {
		r := s.RabbitMQClient(ctx)
		s.eventConsumer = rmqConsumer.NewUserConsumer(r.Connect().Channel, s.EventService(ctx), "post")
	}
	return s.eventConsumer
}
