package routes

import (
	"Ecadr/internal/controllers"
	"Ecadr/internal/controllers/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) *gin.Engine {
	pingRoute := r.Group("/ping")
	{
		pingRoute.GET("/", controllers.Ping)
	}

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
	r.GET("/vacancy", controllers.GetAllVacancies)
	r.GET("/vacancy/:id", controllers.GetVacancyByID)
	vacancy := r.Group("/vacancy", middlewares.CheckUserWorker)
	{
		vacancy.POST("", controllers.CreateVacancy)
		vacancy.PATCH("/:id", controllers.UpdateVacancy)
		vacancy.DELETE("/:id", controllers.DeleteVacancyByID)
	}

	// course Маршруты для курсов
	r.GET("/course", controllers.GetAllCourses)
	r.GET("/course/:id", controllers.GetCourseByID)

	course := r.Group("/course", middlewares.CheckUserWorker)
	{
		course.POST("", controllers.CreateCourse)
		course.PATCH("/:id", controllers.UpdateCourse)
		course.DELETE("/:id", controllers.DeleteCourse)
	}

	return r
}
