package main

import (
	"context"
	"log"
	"os"

	"github.com/ericolvr/sec-back-v2/config"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/database"
	httpServer "github.com/ericolvr/sec-back-v2/internal/infrastructure/http"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/routes"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/sms"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/storage"
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
	actionPlanActivityRepo := database.NewActionPlanActivityRepository(db)
	activityMediaRepo := database.NewActivityMediaRepository(db)
	riskCategoryRepo := database.NewRiskCategoryRepository(db)
	formulaRepo := database.NewCalculationFormulaRepository(db)
	reportRepo := database.NewAnalyticsReportRepository(db)
	assignmentRepo := database.NewAssessmentAssignmentRepository(db)
	riskMetricsRepo := database.NewRiskMetricsRepository(db)
	versionRepo := database.NewAssessmentVersionRepository(db)
	invitationRepo := database.NewInvitationRepository(db)

	// Storage Client (Google Cloud Storage)
	ctx := context.Background()
	storageClient, err := storage.NewStorageClient(ctx, cfg.App.GCSBucket)
	if err != nil {
		log.Printf("Warning: Failed to initialize storage client: %v", err)
		log.Printf("Upload de arquivos não estará disponível")
	}
	if storageClient != nil {
		defer storageClient.Close()
	}

	// Services
	smsProvider := sms.NewTwilioProvider()
	answerService := services.NewAnswerService(answerRepo, questionRepo, submissionRepo)
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
	assignmentService := services.NewAssessmentAssignmentService(assignmentRepo, departmentRepo, employeeRepo, submissionRepo, invitationRepo)
	invitationService := services.NewInvitationService(invitationRepo, submissionRepo, employeeRepo)

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

	calculationFormulaService := services.NewCalculationFormulaService(formulaRepo)

	emailService := services.NewEmailService()

	// Handlers
	answerHandler := api.NewAnswerHandler(answerService, submissionService)
	authHandler := api.NewAuthHandler(authService)
	departmentHandler := api.NewDepartmentHandler(departmentService)
	employeeHandler := api.NewEmployeeHandler(employeeService)
	questionHandler := api.NewQuestionHandler(questionService)
	userHandler := api.NewUserHandler(userService)
	analyticsHandler := api.NewAnalyticsHandler(analyticsService)
	actionPlanHandler := api.NewActionPlanHandler(actionPlanService)
	actionPlanTemplateHandler := api.NewActionPlanTemplateHandler(actionPlanTemplateRepo)
	actionPlanActivityHandler := api.NewActionPlanActivityHandler(actionPlanActivityRepo)
	activityMediaHandler := api.NewActivityMediaHandler(activityMediaRepo, storageClient)
	companyHandler := api.NewCompanyHandler(companyService)
	partnerHandler := api.NewPartnerHandler(partnerService)
	templateHandler := api.NewAssessmentTemplateHandler(templateService)
	versionHandler := api.NewAssessmentVersionHandler(versionService)
	submissionHandler := api.NewEmployeeSubmissionHandler(submissionService)
	assignmentHandler := api.NewAssessmentAssignmentHandler(assignmentService)
	invitationHandler := api.NewInvitationHandler(invitationService, assignmentService, submissionService, emailService)
	emailHandler := api.NewEmailHandler(emailService)
	surveyHandler := api.NewSurveyHandler(submissionService, questionService, templateService)
	calculationFormulaHandler := api.NewCalculationFormulaHandler(calculationFormulaService)

	// Routes
	router := httpServer.NewRouter()
	router.AnswerRoutes = routes.NewAnswerRoutes(answerHandler)
	router.AuthRoutes = routes.NewAuthRoutes(authHandler)
	router.DepartmentRoutes = routes.NewDepartmentRoutes(departmentHandler)
	router.EmployeeRoutes = routes.NewEmployeeRoutes(employeeHandler)
	router.QuestionRoutes = routes.NewQuestionRoutes(questionHandler)
	router.UserRoutes = routes.NewUserRoutes(userHandler)
	router.AnalyticsRoutes = routes.NewAnalyticsRoutes(analyticsHandler)
	router.ActionPlanRoutes = routes.NewActionPlanRoutes(actionPlanHandler, actionPlanActivityHandler)
	router.ActionPlanTemplateRoutes = routes.NewActionPlanTemplateRoutes(actionPlanTemplateHandler)
	router.ActionPlanActivityRoutes = routes.NewActionPlanActivityRoutes(actionPlanActivityHandler)
	router.ActivityMediaRoutes = routes.NewActivityMediaRoutes(activityMediaHandler)
	router.CompanyRoutes = routes.NewCompanyRoutes(companyHandler)
	router.PartnerRoutes = routes.NewPartnerRoutes(partnerHandler)
	router.AssessmentTemplateRoutes = routes.NewAssessmentTemplateRoutes(templateHandler, versionHandler, questionHandler)
	router.EmployeeSubmissionRoutes = routes.NewEmployeeSubmissionRoutes(submissionHandler)
	router.AssessmentAssignmentRoutes = routes.NewAssessmentAssignmentRoutes(assignmentHandler)
	router.InvitationRoutes = routes.NewInvitationRoutes(invitationHandler)
	router.EmailRoutes = routes.NewEmailRoutes(emailHandler)
	router.SurveyRoutes = routes.NewSurveyRoutes(surveyHandler)
	router.CalculationFormulaRoutes = routes.NewCalculationFormulaRoutes(calculationFormulaHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	server := httpServer.NewServer(router, port)
	log.Fatal(server.Start())
}
