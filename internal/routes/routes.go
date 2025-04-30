package routes

import (
	_ "Ecadr/docs"
	"Ecadr/internal/controllers"
	"Ecadr/internal/controllers/ai"
	"Ecadr/internal/controllers/middlewares"
	aiWebSocket "Ecadr/internal/controllers/websockets/ai"
	"Ecadr/internal/controllers/websockets/pagination"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// InitRoutes — настраиваем HTTP-маршруты
func InitRoutes(r *gin.Engine) *gin.Engine {
	// Health-check
	r.GET("/ping", controllers.Ping)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Пользователи
	users := r.Group("/users",
		middlewares.CheckUserAuthentication,
		middlewares.CheckUserWorker,
	)
	{
		users.GET("", pagination.UsersWebSocket)
		users.GET("/:id", controllers.GetUserByID)
	}
	r.GET("/user", middlewares.CheckUserAuthentication, controllers.GetMyData)

	// Auth
	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", controllers.SignUp)
		auth.POST("/sign-in", controllers.SignIn)
		auth.POST("/refresh", controllers.RefreshToken)
	}

	// Компания
	company := r.Group("/company")
	{
		company.GET("/:id", controllers.GetCompanyByID)
		company.PATCH("/:id",
			middlewares.CheckUserAuthentication,
			middlewares.CheckUserWorker,
			middlewares.CheckWorkerCompany,
			controllers.UpdateCompany,
		)
	}

	// Вакансии
	vacancy := r.Group("/vacancy")
	{
		vacancy.GET("",
			middlewares.CheckUserAuthentication,
			func(c *gin.Context) {
				runHeavyTask(c, ai.GetAnalyseForUserVacancies)
			},
		)
		vacancy.GET("/:id", controllers.GetVacancyByID)
		vacancy.GET("/company/:company_id", controllers.GetAllCompanyVacancies)
	}

	// Вакансии (worker)
	vacancyWorker := r.Group("/vacancy",
		middlewares.CheckUserAuthentication,
		middlewares.CheckUserWorker,
		middlewares.CheckWorkerVacancy,
	)
	{
		vacancyWorker.POST("", controllers.CreateVacancy)
		vacancyWorker.PATCH("/:id", controllers.UpdateVacancy)
		vacancyWorker.DELETE("/:id", controllers.DeleteVacancyByID)
		vacancyWorker.GET("/users/:id", ai.GetAnalyseForVacanciesUser)
		vacancyWorker.POST("/recommends", controllers.CreateRecommendVacancy)
	}

	// Критерии вакансий
	criteria := vacancy.Group("/criteria")
	{
		criteria.GET("/:id", controllers.GetVacancyCriteria)
		criteria.GET("/spec/:id", controllers.GetVacancyCriteriaByID)
		criteria.POST("",
			middlewares.CheckUserAuthentication,
			middlewares.CheckUserWorker,
			controllers.CreateVacancyCriteria,
		)
		criteria.PATCH("/:id",
			middlewares.CheckUserAuthentication,
			middlewares.CheckUserWorker,
			middlewares.CheckWorkerVacancyCriteria,
			controllers.UpdateVacancyCriteria,
		)
		criteria.DELETE("/:id",
			middlewares.CheckUserAuthentication,
			middlewares.CheckUserWorker,
			middlewares.CheckWorkerVacancyCriteria,
			controllers.DeleteVacancyCriteria,
		)
	}

	// Курсы
	course := r.Group("/course")
	{
		course.GET("",
			middlewares.CheckUserAuthentication,
			func(c *gin.Context) {
				runHeavyTask(c, ai.GetAnalyseForUserCourse)
			},
		)
		course.GET("/:id", controllers.GetCourseByID)
		course.GET("/worker/:worker_id", controllers.GetAllWorkerCourses)
	}

	// Курсы (worker)
	courseWorker := r.Group("/course",
		middlewares.CheckUserAuthentication,
		middlewares.CheckUserWorker,
	)
	{
		courseWorker.POST("", controllers.CreateCourse)
		courseWorker.PATCH("/:id", middlewares.CheckWorkerCourse, controllers.UpdateCourse)
		courseWorker.DELETE("/:id", middlewares.CheckWorkerCourse, controllers.DeleteCourse)
		courseWorker.GET("/users/:id", middlewares.CheckWorkerCourse, ai.GetAnalyseForCourseUser)
		courseWorker.POST("/recommends", controllers.CreateRecommendCourse)
	}

	// Рекомендации
	recommend := r.Group("/recommends", middlewares.CheckUserAuthentication)
	{
		recommend.GET("/:id", controllers.GetUserRecommendByID)
		recommend.GET("/course", controllers.GetUserRecommendsCourse)
		recommend.GET("/vacancy", controllers.GetUserRecommendsVacancy)
	}

	// AI-Chat
	aiChat := r.Group("/ai", middlewares.CheckUserAuthentication)
	{
		aiChat.GET("/recommends", func(c *gin.Context) {
			runHeavyTask(c, ai.GetAIRecommendsForUser)
		})
		aiChat.GET("/chat", aiWebSocket.ChatAIWebSocketHandler)
	}

	// Аналитика
	analyse := r.Group("/analyse",
		middlewares.CheckUserAuthentication,
		middlewares.CheckWorkerDepartment,
	)
	{
		analyse.GET("/companies", func(c *gin.Context) {
			runHeavyTask(c, ai.GetCompaniesStatistic)
		})
		analyse.GET("/users", func(c *gin.Context) {
			runHeavyTask(c, ai.GetUsersStatistic)
		})
	}

	// Тесты
	tests := r.Group("/tests")
	{
		tests.POST("", controllers.CreateTest)
		tests.GET("", controllers.GetTestsByTypeAndID)
		tests.GET("/:id", controllers.GetTestByID)

		tests.POST("", middlewares.CheckUserAuthentication, middlewares.CheckUserWorker, middlewares.CheckUserTest, controllers.CreateTest)
		tests.PUT("/:id", middlewares.CheckUserAuthentication, middlewares.CheckUserWorker, middlewares.CheckUserTest, controllers.UpdateTest)
		tests.DELETE("/:id", middlewares.CheckUserAuthentication, middlewares.CheckUserWorker, middlewares.CheckUserTest, controllers.DeleteTest)
	}

	testSubmissions := tests.Group("/submissions")
	{
		testSubmissions.GET("/:id", controllers.GetSortedSubmissionsHandler)
		testSubmissions.POST("", middlewares.CheckUserAuthentication, controllers.CreateTestSubmission)
	}

	return r
}
