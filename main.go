package main

import (
	"log"
	"os"

	_ "github.com/hisyntax/agritech/docs"
	"github.com/hisyntax/agritech/user"
	"github.com/hisyntax/monnify-go"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           agritech API
// @version         1.0
// @description     This is the API serving the agritech frontend
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host                       agritech0.herokuapp.com
// @BasePath                   /api/v1
// @schemes                    https
// @query.collection.format    multi
// @securityDefinitions.basic  BasicAuth
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}

	apiKey := os.Getenv("MONNIFY_API_KEY")
	sKey := os.Getenv("MONNIFY_SECRET_KEY")
	bUrl := os.Getenv("MONNIFY_BASE_URL")
	monnify.Options(apiKey, sKey, bUrl)
}

func main() {
	r := gin.Default()
	config := CORSMiddleware()
	r.Use(config)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	var (
		users = user.User{}
	)

	//api version
	api := r.Group("/api/v1")
	// api.
	{
		user := api.Group("/user")
		{
			user.POST("/signup", users.Signup)
			user.POST("/signin", users.Signin)
			user.POST("/user/otp", users.SendOTP)
			user.POST("/user/ot/validatep", users.ValidateOtp)
		}
	}

	r.GET("/api/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":" + port)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding , X-CSRF-Token, Authorization, redirect_url ,email, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
