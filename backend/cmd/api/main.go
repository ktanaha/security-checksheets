package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/security-checksheets/backend/internal/infrastructure/repository"
	"github.com/security-checksheets/backend/internal/interface/handler"
	"github.com/security-checksheets/backend/internal/usecase"
)

func main() {
	// Ginモードの設定
	if os.Getenv("GO_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// データベース接続
	db := initDB()
	defer db.Close()

	// 依存性の注入（Clean Architecture）
	projectRepo := repository.NewProjectRepository(db)
	projectUseCase := usecase.NewProjectUseCase(projectRepo)
	projectHandler := handler.NewProjectHandler(projectUseCase)

	// ファイル管理
	fileRepo := repository.NewFileRepository(db)
	fileUseCase := usecase.NewFileUseCase(fileRepo, projectRepo, "/app/uploads")
	fileHandler := handler.NewFileHandler(fileUseCase)

	// Ginルーターの初期化
	router := gin.Default()

	// CORS設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ヘルスチェックエンドポイント
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "security-checksheets-backend",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// APIルートグループ
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// 案件管理エンドポイント
		projects := api.Group("/projects")
		{
			projects.POST("", projectHandler.CreateProject)
			projects.GET("", projectHandler.ListProjects)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)

			// 案件に紐づくファイル管理
			projects.POST("/:id/files", fileHandler.UploadFile)
			projects.GET("/:id/files", fileHandler.ListFilesByProject)
		}

		// ファイル管理エンドポイント
		files := api.Group("/files")
		{
			files.GET("/:id", fileHandler.GetFile)
			files.DELETE("/:id", fileHandler.DeleteFile)
		}
	}

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("サーバーをポート %s で起動しています...", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}

// initDB はデータベース接続を初期化する
func initDB() *sql.DB {
	// 環境変数からDB接続情報を取得
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "admin"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "security_checksheets"
	}

	// 接続文字列の構築
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// データベース接続
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("データベースへのPingに失敗しました: %v", err)
	}

	log.Println("データベース接続に成功しました")
	return db
}
