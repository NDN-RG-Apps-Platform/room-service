package config

import (
	"os"
	"room-service/common/util"

	"github.com/sirupsen/logrus"
)

var Config AppConfig

type AppConfig struct {
	Port                  int             `json:"port"`
	AppName               string          `json:"appName"`
	AppEnv                string          `json:"appEnv"`
	SignatureKey          string          `json:"signatureKey"`
	Database              Database        `json:"database"`
	RateLimiterMaxRequest float64         `json:"rateLimiterMaxRequest"` // Dikembalikan ke float64
	RateLimiterTimeSecond int             `json:"rateLimiterTimeSecond"`
	InternalService       InternalService `json:"internalService"`
	GcsType               string          `json:"gcsType"`
	GcsProjectID          string          `json:"gcsProjectID"`
	GcsPrivateKeyID       string          `json:"gcsPrivateKeyID"`
	GcsPrivateKey         string          `json:"gcsPrivateKey"`
	GcsClientEmail        string          `json:"gcsClientEmail"`
	GcsClientID           string          `json:"gcsClientID"`
}

type InternalService struct {
	User struct {
		Host         string `json:"host"`
		SignatureKey string `json:"signatureKey"`
	} `json:"user"`
}

type Database struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	Name                  string `json:"name"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	MaxOpenConnection     int    `json:"maxOpenConnection"`
	MaxLifeTimeConnection int    `json:"maxLifeTimeConnection"`
	MaxIdleConnection     int    `json:"maxIdleConnection"`
	MaxIdleTime           int    `json:"maxIdleTime"`
}

func Init() {
	err := util.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		logrus.Errorf("Failed to bind config from JSON: %v", err) // Gunakan logrus.Error untuk log error
		err = util.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			logrus.Fatalf("Failed to bind config from Consul: %v", err) // Fatalf menghentikan program
		}
	}
}
