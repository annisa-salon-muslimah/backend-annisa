package handler

import (
	"annisa-salon/auth"
	"annisa-salon/database"
	"annisa-salon/middleware"
	"annisa-salon/repository"
	"annisa-salon/service"
	"log"
	"os"

	// _ "annisa-salon/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartApp() {
	db, err := database.InitDb()
	if err != nil {
		log.Fatal("Eror Db Connection")
	}

	secretKey := os.Getenv("SECRET_KEY")

	router := gin.Default()

		//add sweager
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Origin , Accept , X-Requested-With , Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization"},
		AllowMethods:    []string{"POST, OPTIONS, GET, PUT, DELETE"},
	}))
	
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := auth.NewUserAuthService()
	authService.SetSecretKey(secretKey)
	userHandler := NewUserHandler(userService, authService)

	blogRepository := repository.NewBlogRepository(db)
	blogService := service.NewBlogService(blogRepository)
	blogHandler := NewBlogHandler(blogService, authService)

	treatmentRepository := repository.NewTreatmentsRepository(db)
	treatmentService := service.NewTreatmentService(treatmentRepository)
	treatmentHandler := NewTreatmentsHandler(treatmentService, authService)

	user := router.Group("api/user")
	user.POST("/register", userHandler.RegisterUser)
	user.POST("/login", userHandler.Login)

	blog := router.Group("api/blog")
	blog.POST("/", middleware.AuthMiddleware(authService, userService), middleware.AuthRole(authService, userService), blogHandler.CreateBlog)
	blog.PUT("/:slug", middleware.AuthMiddleware(authService, userService), middleware.AuthRole(authService, userService),blogHandler.UpdateBlog)
	blog.GET("/:slug", blogHandler.GetOneBlog)
	blog.GET("/", blogHandler.GetAllBlog)
	blog.DELETE("/:slug", middleware.AuthMiddleware(authService, userService), middleware.AuthRole(authService, userService), blogHandler.DeleteBlog)

	treatment := router.Group("api/treatment")
	treatment.POST("/", middleware.AuthMiddleware(authService, userService), middleware.AuthRole(authService, userService), treatmentHandler.CreateTreatments)
	treatment.PUT("/:slug", middleware.AuthMiddleware(authService, userService), middleware.AuthRole(authService, userService), treatmentHandler.UpdatedTreatment)
	treatment.GET("/:slug", treatmentHandler.GetOneTreatment)
	treatment.GET("/", treatmentHandler.GetAllTreatments)
	treatment.DELETE("/:slug", middleware.AuthMiddleware(authService, userService), middleware.AuthRole(authService, userService), treatmentHandler.DeleteTreatment)

	router.Run(":8080")

}