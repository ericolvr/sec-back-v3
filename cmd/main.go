package main

import (
	"log"
	"os"

	"github.com/ericolvr/sec-back-v2/config"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/database"
	httpServer "github.com/ericolvr/sec-back-v2/internal/infrastructure/http"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/routes"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/sms"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	db, err := database.Connect(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	answerRepo := database.NewAnswerRepository(db)
	departmentRepo := database.NewDepartmentRepository(db)
	employeeRepo := database.NewEmployeeRepository(db)
	questionRepo := database.NewQuestionRepository(db)
	userRepo := database.NewUserRepository(db)

	smsProvider := sms.NewTwilioProvider()
	answerService := services.NewAnswerService(answerRepo, questionRepo, nil)
	authService := services.NewAuthService(userRepo)
	departmentService := services.NewDepartmentService(departmentRepo)
	employeeService := services.NewEmployeeService(employeeRepo)
	questionService := services.NewQuestionService(questionRepo)
	userService := services.NewUserService(userRepo, smsProvider)

	answerHandler := api.NewAnswerHandler(answerService, nil)
	authHandler := api.NewAuthHandler(authService)
	departmentHandler := api.NewDepartmentHandler(departmentService)
	employeeHandler := api.NewEmployeeHandler(employeeService)
	questionHandler := api.NewQuestionHandler(questionService)
	userHandler := api.NewUserHandler(userService)

	router := httpServer.NewRouter()
	router.AnswerRoutes = routes.NewAnswerRoutes(answerHandler)
	router.AuthRoutes = routes.NewAuthRoutes(authHandler)
	router.DepartmentRoutes = routes.NewDepartmentRoutes(departmentHandler)
	router.EmployeeRoutes = routes.NewEmployeeRoutes(employeeHandler)
	router.QuestionRoutes = routes.NewQuestionRoutes(questionHandler)
	router.UserRoutes = routes.NewUserRoutes(userHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	server := httpServer.NewServer(router, port)
	log.Fatal(server.Start())
}
