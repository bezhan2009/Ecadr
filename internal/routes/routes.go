package routes

import (
	_ "Ecadr/docs"
	"Ecadr/internal/controllers"
	"Ecadr/internal/controllers/ai"
	"Ecadr/internal/controllers/middlewares"
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

	// vacancy Маршруты для вакансий
	r.GET("/vacancy/:id", controllers.GetVacancyByID)
	r.GET("/vacancy", middlewares.CheckUserAuthentication, ai.GetAnalyseForUserVacancies)
	vacancy := r.Group("/vacancy", middlewares.CheckUserWorker)
	{
		vacancy.POST("", controllers.CreateVacancy)
		vacancy.PATCH("/:id", controllers.UpdateVacancy)
		vacancy.DELETE("/:id", controllers.DeleteVacancyByID)
	}

	// course Маршруты для курсов
	r.GET("/course", middlewares.CheckUserAuthentication, ai.GetAnalyseForUserCourse)
	r.GET("/course/:id", controllers.GetCourseByID)
	course := r.Group("/course", middlewares.CheckUserWorker)
	{
		course.POST("", controllers.CreateCourse)
		course.PATCH("/:id", controllers.UpdateCourse)
		course.DELETE("/:id", controllers.DeleteCourse)
	}

	return r
}
