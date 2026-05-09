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

type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB
	Redis  *redis.Client
	Store  *minio.Client
	Router *gin.Engine
}

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
	if err := ensureSchema(db); err != nil {
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

func ensureSchema(db *gorm.DB) error {
	if err := db.Exec("ALTER TABLE sys_users ADD COLUMN IF NOT EXISTS avatar_id BIGINT").Error; err != nil {
		return err
	}
	if err := ensureFileDeletePermission(db); err != nil {
		return err
	}
	return db.Exec(`
DO $$
DECLARE
  item RECORD;
BEGIN
  FOR item IN
    SELECT table_name, column_name
    FROM (VALUES
      ('sys_depts', 'created_at'),
      ('sys_depts', 'updated_at'),
      ('sys_users', 'created_at'),
      ('sys_users', 'updated_at'),
      ('sys_roles', 'created_at'),
      ('sys_roles', 'updated_at'),
      ('sys_menus', 'created_at'),
      ('sys_menus', 'updated_at'),
      ('sys_files', 'created_at'),
      ('sys_files', 'updated_at'),
      ('sys_login_logs', 'created_at'),
      ('sys_operation_logs', 'created_at')
    ) AS columns(table_name, column_name)
  LOOP
    IF EXISTS (
      SELECT 1
      FROM information_schema.columns
      WHERE table_schema = 'public'
        AND table_name = item.table_name
        AND column_name = item.column_name
        AND data_type = 'bigint'
    ) THEN
      EXECUTE format('ALTER TABLE %I ALTER COLUMN %I DROP DEFAULT', item.table_name, item.column_name);
      EXECUTE format('ALTER TABLE %I ALTER COLUMN %I TYPE TIMESTAMPTZ USING to_timestamp(%I / 1000.0)', item.table_name, item.column_name, item.column_name);
      EXECUTE format('ALTER TABLE %I ALTER COLUMN %I SET DEFAULT CURRENT_TIMESTAMP', item.table_name, item.column_name);
    END IF;
  END LOOP;
END $$;
`).Error
}

func ensureFileDeletePermission(db *gorm.DB) error {
	if err := db.Exec(`
INSERT INTO sys_menus(id, parent_id, name, title, type, path, component, icon, permission, sort)
VALUES (122, 7, 'FileDelete', '删除文件', 'button', '', '', '', 'file:delete', 2)
ON CONFLICT (id) DO NOTHING
`).Error; err != nil {
		return err
	}
	return db.Exec(`
INSERT INTO sys_role_menus(role_id, menu_id)
SELECT 1, 122
WHERE EXISTS (SELECT 1 FROM sys_roles WHERE id = 1)
ON CONFLICT DO NOTHING
`).Error
}

func (a *App) Run() error {
	srv := &http.Server{
		Addr:              a.Config.Server.Addr,
		Handler:           a.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	a.Logger.Info("server started", zap.String("addr", a.Config.Server.Addr))
	return srv.ListenAndServe()
}

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
	authHandler.RegisterPrivate(private)
	systemHandler.Register(private, require)
	filemod.NewHandler(a.DB, a.Store, a.Config.RustFS).Register(private, require)
	audit.NewHandler(a.DB).Register(private, require)

	return r
}

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
