package cmd

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"room-service/clients"
	"room-service/common/gcs"
	"room-service/common/response"
	"room-service/config"
	"room-service/constants"
	controllers "room-service/controllers"
	"room-service/domain/models"
	"room-service/middlewares"
	"room-service/repositories"
	"room-service/routes"
	"room-service/services"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var commad = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(c *cobra.Command, args []string) {
		_ = godotenv.Load()
		config.Init()
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		err = db.AutoMigrate(
			&models.Room{},
			&models.RoomSchedule{},
			&models.Time{},
		)
		if err != nil {
			panic(err)
		}

		gcs := initGCS()
		client := clients.NewClientRegistry()
		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository, gcs)
		controller := controllers.NewControllerRegistry(service)

		router := gin.Default()
		router.Use(middlewares.HandlePanic())
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})

		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to user service",
			})
		})

		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-api-key, x-request-at")
			c.Next()
		})

		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			})

		router.Use(middlewares.RateLimiter(lmt))

		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group, client)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)
	},
}

func Run() {
	err := commad.Execute()
	if err != nil {
		panic(err)
	}
}

func initGCS() gcs.IGCSClient {
	decode, err := base64.StdEncoding.DecodeString(config.Config.GcsPrivateKey)
	if err != nil {
		panic(err)
	}
	stringPriviteKey := string(decode)
	gcsServiceAccount := gcs.ServiceAccountKeyJSON{
		Type:                    config.Config.GcsType,
		ProjectID:               config.Config.GcsProjectID,
		PrivateKeyID:            config.Config.GcsPrivateKeyID,
		PrivateKey:              stringPriviteKey,
		ClientEmail:             config.Config.GcsClientEmail,
		ClientID:                config.Config.GcsClientID,
		AuthURI:                 config.Config.GcsAuthUri,
		TokenURI:                config.Config.GcsTokenURI,
		AuthProviderX509CertUrl: config.Config.GcsAuthProviderX509CertUrl,
		ClientX509CertUrl:       config.Config.GcsClientX509CertUrl,
		UniverseDomain:          config.Config.GcsUniverseDomain,
	}
	gcsClient := gcs.NewGCSClient(
		gcsServiceAccount,
		config.Config.gcsBucketName,
	)
	return gcsClient
}
