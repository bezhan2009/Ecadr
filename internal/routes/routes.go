package routes

import (
	_ "Ecadr/docs"
	"Ecadr/internal/controllers"
	"Ecadr/internal/controllers/ai"
	"Ecadr/internal/controllers/middlewares"
	aiWebSocket "Ecadr/internal/controllers/websockets/ai"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func InitRoutes(r *gin.Engine) *gin.Engine {
	pingRoute := r.Group("/ping")
	{
		pingRoute.GET("/", controllers.Ping)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// usersRoute Маршруты для пользователей (профили)
	usersRoute := r.Group("/users")
	{
		usersRoute.GET("", controllers.GetAllUsers)
		usersRoute.GET("/:id", controllers.GetUserByID)
	}

	r.GET("user", middlewares.CheckUserAuthentication, controllers.GetMyData)

	// auth Маршруты для авторизаций
	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", controllers.SignUp)
		auth.POST("/sign-in", controllers.SignIn)
		auth.POST("/refresh", controllers.RefreshToken)
	}

	// vacancy и vacancyWorker Маршруты для вакансий
	vacancy := r.Group("/vacancy")
	{
		vacancy.GET("", middlewares.CheckUserAuthentication,
			ai.GetAnalyseForUserVacancies)
		vacancy.GET("/:id", controllers.GetVacancyByID)
		vacancy.GET("/worker/:worker_id", controllers.GetAllWorkerVacancies)
	}

	vacancyWorker := r.Group("/vacancy", middlewares.CheckUserAuthentication, middlewares.CheckUserWorker)
	{
		vacancyWorker.POST("", controllers.CreateVacancy)
		vacancyWorker.PATCH("/:id", controllers.UpdateVacancy)
		vacancyWorker.DELETE("/:id", controllers.DeleteVacancyByID)

		vacancyWorker.GET("/users/:id", ai.GetAnalyseForVacanciesUser)
		vacancyWorker.POST("/recommends", controllers.CreateRecommendVacancy)
	}

	// course и courseWorker Маршруты для курсов
	course := r.Group("/course")
	{
		course.GET("", middlewares.CheckUserAuthentication, ai.GetAnalyseForUserCourse)
		course.GET("/:id", controllers.GetCourseByID)
		course.GET("/worker/:worker_id", controllers.GetAllWorkerCourses)
	}

	courseWorker := r.Group("/course", middlewares.CheckUserAuthentication, middlewares.CheckUserWorker)
	{
		courseWorker.POST("", controllers.CreateCourse)
		courseWorker.PATCH("/:id", controllers.UpdateCourse)
		courseWorker.DELETE("/:id", controllers.DeleteCourse)

		courseWorker.POST("/recommends", controllers.CreateRecommendCourse)
	}

	// recommends Маршруты для получения своих рекомендаций
	recommend := r.Group("/recommends", middlewares.CheckUserAuthentication)
	{
		recommend.GET("/:id", controllers.GetUserRecommendByID)
		recommend.GET("/course", controllers.GetUserRecommendsCourse)
		recommend.GET("/vacancy", controllers.GetUserRecommendsVacancy)
	}

	// aiChat Маршруты для ИИ
	aiChat := r.Group("ai", middlewares.CheckUserAuthentication)
	{
		aiChat.GET("/recommends", ai.GetAIRecommendsForUser)

		aiChat.GET("/chat", aiWebSocket.ChatAIWebSocketHandler)
	}

	return r
}
