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

	// Repositories
	answerRepo := database.NewAnswerRepository(db)
	departmentRepo := database.NewDepartmentRepository(db)
	employeeRepo := database.NewEmployeeRepository(db)
	questionRepo := database.NewQuestionRepository(db)
	userRepo := database.NewUserRepository(db)
	companyRepo := database.NewCompanyRepository(db)
	partnerRepo := database.NewPartnerRepository(db)
	submissionRepo := database.NewEmployeeSubmissionRepository(db)
	templateRepo := database.NewAssessmentTemplateRepository(db)
	settingsRepo := database.NewPartnerSettingsRepository(db)
	actionPlanRepo := database.NewActionPlanRepository(db)
	actionPlanTemplateRepo := database.NewActionPlanTemplateRepository(db)
	riskCategoryRepo := database.NewRiskCategoryRepository(db)
	formulaRepo := database.NewCalculationFormulaRepository(db)
	reportRepo := database.NewAnalyticsReportRepository(db)
	assignmentRepo := database.NewQuestionnaireAssignmentRepository(db)
	riskMetricsRepo := database.NewRiskMetricsRepository(db)
	versionRepo := database.NewAssessmentVersionRepository(db)

	// Services
	smsProvider := sms.NewTwilioProvider()
	answerService := services.NewAnswerService(answerRepo, questionRepo, nil)
	authService := services.NewAuthService(userRepo)
	departmentService := services.NewDepartmentService(departmentRepo)
	employeeService := services.NewEmployeeService(employeeRepo)
	questionService := services.NewQuestionService(questionRepo)
	userService := services.NewUserService(userRepo, smsProvider)
	companyService := services.NewCompanyService(companyRepo, partnerRepo)
	partnerService := services.NewPartnerService(partnerRepo)
	versionService := services.NewAssessmentVersionService(versionRepo)
	templateService := services.NewAssessmentTemplateService(templateRepo, partnerRepo, versionService)
	submissionService := services.NewEmployeeSubmissionService(submissionRepo, employeeRepo, templateRepo)
	assignmentService := services.NewQuestionnaireAssignmentService(assignmentRepo, departmentRepo)

	// RiskMetricsService e AnalyticsService
	riskMetricsService := services.NewRiskMetricsService(
		riskMetricsRepo,
		submissionRepo,
		answerRepo,
		questionRepo,
		employeeRepo,
		departmentRepo,
		riskCategoryRepo,
		formulaRepo,
		nil, // actionPlanService - será nil por enquanto
	)

	analyticsService := services.NewAnalyticsService(
		riskMetricsService,
		departmentRepo,
		companyRepo,
		partnerRepo,
		templateRepo,
		settingsRepo,
		actionPlanRepo,
		actionPlanTemplateRepo,
		riskCategoryRepo,
		questionRepo,
		submissionRepo,
		answerRepo,
		formulaRepo,
		reportRepo,
		assignmentRepo,
	)

	actionPlanService := services.NewActionPlanService(
		actionPlanRepo,
		actionPlanTemplateRepo,
		riskCategoryRepo,
		departmentRepo,
	)

	// Handlers
	answerHandler := api.NewAnswerHandler(answerService, nil)
	authHandler := api.NewAuthHandler(authService)
	departmentHandler := api.NewDepartmentHandler(departmentService)
	employeeHandler := api.NewEmployeeHandler(employeeService)
	questionHandler := api.NewQuestionHandler(questionService)
	userHandler := api.NewUserHandler(userService)
	analyticsHandler := api.NewAnalyticsHandler(analyticsService)
	actionPlanHandler := api.NewActionPlanHandler(actionPlanService)
	companyHandler := api.NewCompanyHandler(companyService)
	partnerHandler := api.NewPartnerHandler(partnerService)
	templateHandler := api.NewAssessmentTemplateHandler(templateService)
	versionHandler := api.NewAssessmentVersionHandler(versionService)
	submissionHandler := api.NewEmployeeSubmissionHandler(submissionService)
	assignmentHandler := api.NewQuestionnaireAssignmentHandler(assignmentService)

	// Routes
	router := httpServer.NewRouter()
	router.AnswerRoutes = routes.NewAnswerRoutes(answerHandler)
	router.AuthRoutes = routes.NewAuthRoutes(authHandler)
	router.DepartmentRoutes = routes.NewDepartmentRoutes(departmentHandler)
	router.EmployeeRoutes = routes.NewEmployeeRoutes(employeeHandler)
	router.QuestionRoutes = routes.NewQuestionRoutes(questionHandler)
	router.UserRoutes = routes.NewUserRoutes(userHandler)
	router.AnalyticsRoutes = routes.NewAnalyticsRoutes(analyticsHandler)
	router.ActionPlanRoutes = routes.NewActionPlanRoutes(actionPlanHandler)
	router.CompanyRoutes = routes.NewCompanyRoutes(companyHandler)
	router.PartnerRoutes = routes.NewPartnerRoutes(partnerHandler)
	router.AssessmentTemplateRoutes = routes.NewAssessmentTemplateRoutes(templateHandler, versionHandler, questionHandler)
	router.EmployeeSubmissionRoutes = routes.NewEmployeeSubmissionRoutes(submissionHandler)
	router.QuestionnaireAssignmentRoutes = routes.NewQuestionnaireAssignmentRoutes(assignmentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	server := httpServer.NewServer(router, port)
	log.Fatal(server.Start())
}
