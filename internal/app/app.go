package app

import (
	"fmt"
	"github.com/Dor1ma/Time-Tracker/config"
	"github.com/Dor1ma/Time-Tracker/internal/handlers"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"github.com/Dor1ma/Time-Tracker/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Start() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbName,
		cfg.DbPort)

	var db *gorm.DB

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	RunMigration(dsn)

	userRepository := repositories.NewUserRepository(db)
	taskRepository := repositories.NewTaskRepositoryImpl(db)

	userService := services.NewUserServiceImpl(userRepository, cfg.ExternalAPIURL)
	taskService := services.NewTaskServiceImpl(taskRepository)

	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)

	router := gin.Default()

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
		log.Fatalf("Error in gin run function", err.Error())
		return
	}

	log.Println("Gin server has started")
}
