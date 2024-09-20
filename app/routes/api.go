package routes

import (
	"github.com/gin-gonic/gin"
	api "github.com/surajNirala/user_services/app/controllers/API"
)

func ApiRoutes(apiRouter *gin.Engine) {
	// apiRouter.GET("/", api.UserList)
	route := apiRouter.Group("/api")
	{
		// route.GET("/", api.UserList)
		route.GET("/users", api.UserList)
		route.POST("/users/store", api.UserStore)
		route.GET("/users/:user_id", api.UserDetail)
		route.PUT("/users/:user_id", api.UserUpdate)
		route.DELETE("/users/:user_id", api.UserDelete)
	}
}
