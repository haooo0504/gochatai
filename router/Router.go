package router

import (
	"gochatai/docs"
	"gochatai/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Static("/assets/images", "./assets/images")
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/ping", service.GetIndex)

	// user
	r.GET("/user/getUserList", service.GetUserList)
	r.GET("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	r.GET("/user/createUser", service.CreateUser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)

	// post
	r.GET("/post/getPostList", service.GetPostList)
	r.POST("/post/createPost", service.CreatePost)

	return r
}
