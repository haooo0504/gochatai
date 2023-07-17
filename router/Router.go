package router

import (
	"fmt"
	"gochatai/docs"
	"gochatai/service"

	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Static("/assets/images", "./assets/images")
	docs.SwaggerInfo.BasePath = ""
	// Public routes
	public := r.Group("/")
	public.GET("/ping", service.GetIndex)
	public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	public.GET("/user/createUser", service.CreateUser)
	public.GET("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	// Private (authenticated) routes
	private := r.Group("/")
	private.Use(JWTAuthMiddleware())
	private.GET("/user/getUserList", service.GetUserList)
	private.GET("/user/deleteUser", service.DeleteUser)
	private.POST("/user/updateUser", service.UpdateUser)

	private.GET("/post/getPostList", service.GetPostList)
	private.POST("/post/createPost", service.CreatePost)
	return r
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, BearerSchema) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerSchema)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Always check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key (this should be read from environment variable or config file)
			// This is just an example and should not be used in production
			return []byte("your_secret_key"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// You can access the token claims (payload) here
			c.Set("userID", claims["userID"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
	}
}
