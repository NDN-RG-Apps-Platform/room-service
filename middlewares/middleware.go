package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"room-service/clients"
	"room-service/common/response"
	"room-service/config"
	"room-service/constants"
	errConstant "room-service/constants/error"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstant.ErrInternalServerError.Error(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimiter(let *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(let, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// func extractBearerToken(token string) string {
// 	arrayToken := strings.Split(token, " ")
// 	if len(arrayToken) == 2 {
// 		return arrayToken[1]
// 	}
// 	return ""
// }

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	c.Abort()
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func CheckRole(roles []string, client clients.IClientRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := client.GetUser().GetUserByToken(c.Request.Context())
		if err != nil {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		if !contains(roles, user.Role) {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}
		c.Next()
	}
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	validateKey := fmt.Sprintf("%s%s%s", serviceName, signatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		return errConstant.ErrUnauthorized
	}
	return nil
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		err = validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}

func AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := validateAPIKey(c)
		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}
