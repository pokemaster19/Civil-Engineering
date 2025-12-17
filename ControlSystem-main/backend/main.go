package main

import (
	"ControlSystem/handlers"
	"ControlSystem/models"
	"ControlSystem/routers"
	"ControlSystem/storage"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v2"
	"gorm.io/gorm"
)

func DbInit() *gorm.DB {
	db, err := models.Setup()
	if err != nil {
		log.Println("Connection error")
	}
	return db
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	db := DbInit()
	minioClient := storage.NewMinioClient(os.Getenv("MINIO_SERVER_URL"), os.Getenv("MINIO_PORT"), os.Getenv("MINIO_ROOT_USER"), os.Getenv("MINIO_ROOT_PASSWORD"), []string{"images", "files"}, false)
	resendClient := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	server := handlers.NewServer(db, minioClient, resendClient)

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			os.Getenv("CLIENT_URL"),
			"http://localhost",
			"http://localhost:80",
			"http://127.0.0.1",
			"http://127.0.0.1:80",
		},
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("api/v1")
	auth := api.Group("/auth")
	routers.RegisterAuthRoutes(auth, server)

	projects := api.Group("/projects")
	routers.RegisterProjectsRoutes(projects, server)

	defects := api.Group("/defects")
	routers.RegisterDefectsRoutes(defects, server)

	attachments := api.Group("/attachments")
	routers.RegisterAttachmentsRoutes(attachments, server)

	admin := api.Group("/admin")
	routers.RegisterAdminRoutes(admin, server)

	return r
}

func main() {
	r := SetupRouter()

	r.Run(":8080")
}
