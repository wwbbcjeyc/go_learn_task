package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/task/go_learn_task/blog-backend/config"
	"github.com/task/go_learn_task/blog-backend/controllers"
	"github.com/task/go_learn_task/blog-backend/database"
	"github.com/task/go_learn_task/blog-backend/middleware"
	"github.com/task/go_learn_task/blog-backend/utils"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 连接数据库
	if err := database.ConnectDB(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// 数据库迁移
	if err := database.MigrateDB(); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 中间件
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())

	// 初始化控制器
	authController := controllers.NewAuthController(cfg)
	postController := controllers.NewPostController()
	commentController := controllers.NewCommentController()

	// 公开路由
	router.POST("/api/register", authController.Register)
	router.POST("/api/login", authController.Login)

	// 文章公开路由
	router.GET("/api/posts", postController.GetAllPosts)
	router.GET("/api/posts/:id", postController.GetPost)

	// 认证路由组
	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleware(cfg))
	{
		// 需要认证的文章操作
		auth.POST("/posts", postController.CreatePost)
		auth.PUT("/posts/:id", postController.UpdatePost)
		auth.DELETE("/posts/:id", postController.DeletePost)

		// 评论操作
		auth.POST("/post-comments/:postId/comments", commentController.CreateComment)
	}

	// 评论公开路由
	router.GET("/api/post-comments/:postId/comments", commentController.GetPostComments)

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, 200, "Server is running", nil)
	})

	// 启动服务器
	log.Printf("🚀 Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
