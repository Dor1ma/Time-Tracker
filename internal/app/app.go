package app

import (
	"fmt"
	"time"

	"github.com/Dor1ma/Time-Tracker/config"
	_ "github.com/Dor1ma/Time-Tracker/docs"
	"github.com/Dor1ma/Time-Tracker/internal/handlers"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"github.com/Dor1ma/Time-Tracker/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func Start() {
	log := logrus.New()
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbName,
		cfg.DbPort)

	var db *gorm.DB
	var retries = 5
	var attempt int

	for attempt = 1; attempt <= retries; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Errorf("Attempt %d failed to connect to database: %v", attempt, err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if attempt > retries {
		log.Fatalf("Failed to connect to database after %d attempts", retries)
		return
	}

	RunMigration(dsn, log)

	userRepository := repositories.NewUserRepositoryImpl(db, log)
	taskRepository := repositories.NewTaskRepositoryImpl(db, log)

	userService := services.NewUserServiceImpl(userRepository, cfg.ExternalAPIURL, log)
	taskService := services.NewTaskServiceImpl(taskRepository, log)

	userHandler := handlers.NewUserHandler(userService, log)
	taskHandler := handlers.NewTaskHandler(taskService, log)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.GET("", userHandler.GetUsers)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.POST("/start", taskHandler.StartTask)
		taskRoutes.POST("/stop", taskHandler.StopTask)
		taskRoutes.GET("/user/:user_id", taskHandler.GetUserTasks)
	}

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error in gin run function: %v", err)
		return
	}

	log.Infof("Gin server has started")
}
