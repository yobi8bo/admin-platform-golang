package bootstrap

import (
	"context"
	"net/http"
	"time"

	"admin-platform/backend/internal/config"
	"admin-platform/backend/internal/middleware"
	"admin-platform/backend/internal/modules/audit"
	"admin-platform/backend/internal/modules/auth"
	filemod "admin-platform/backend/internal/modules/file"
	"admin-platform/backend/internal/modules/system"
	"admin-platform/backend/internal/pkg/logger"
	"admin-platform/backend/internal/pkg/response"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// App 持有后端运行期的共享依赖，路由、数据库、缓存和对象存储都从这里装配。
type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB
	Redis  *redis.Client
	Store  *minio.Client
	Router *gin.Engine
}

// New 根据配置路径初始化应用依赖，并完成私有路由、公共路由和对象存储桶准备。
func New(configPath string) (*App, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}
	log, err := logger.New(cfg.Log)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	store, err := minio.New(cfg.RustFS.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.RustFS.AccessKey, cfg.RustFS.SecretKey, ""),
		Secure: cfg.RustFS.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	if err := ensureBucket(context.Background(), store, cfg.RustFS.Bucket); err != nil {
		return nil, err
	}

	app := &App{Config: cfg, Logger: log, DB: db, Redis: rdb, Store: store}
	app.Router = app.buildRouter()
	return app, nil
}

// Run 启动 HTTP 服务，调用方负责处理 ListenAndServe 返回的错误。
func (a *App) Run() error {
	srv := &http.Server{
		Addr:              a.Config.Server.Addr,
		Handler:           a.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	a.Logger.Info("server started", zap.String("addr", a.Config.Server.Addr))
	return srv.ListenAndServe()
}

// buildRouter 集中注册中间件和模块路由，私有接口统一挂载认证和操作审计。
func (a *App) buildRouter() *gin.Engine {
	if a.Config.Server.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.Trace(), middleware.Recovery(a.Logger), middleware.RequestLogger(a.Logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     a.Config.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/healthz", func(c *gin.Context) { response.OK(c, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	systemHandler := system.NewHandler(a.DB)
	authHandler := auth.NewHandler(a.DB, a.Redis, a.Config.JWT, systemHandler)
	authHandler.RegisterPublic(api)

	private := api.Group("")
	private.Use(middleware.Auth(a.Config.JWT, a.DB), middleware.OperationAudit(a.DB))
	require := func(permission string) gin.HandlerFunc {
		return middleware.RequirePermission(a.DB, permission)
	}
	// 模块路由只接收私有路由组，避免业务模块绕过统一认证和审计边界。
	authHandler.RegisterPrivate(private)
	systemHandler.Register(private, require)
	filemod.NewHandler(a.DB, a.Store, a.Config.RustFS).Register(private, require)
	audit.NewHandler(a.DB).Register(private, require)

	return r
}

// ensureBucket 在启动阶段确保对象存储桶存在，避免首次上传时才暴露存储配置问题。
func ensureBucket(ctx context.Context, client *minio.Client, bucket string) error {
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}
