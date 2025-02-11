package util

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PaginationParam struct {
	Count int64       `json:"count"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}

type PaginationResult struct {
	TotalPage int         `json:"totalPage"`
	TotalData int64       `json:"totalData"`
	NextPage  *int        `json:"nextPage"`
	PrevPage  *int        `json:"prevPage"`
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Data      interface{} `json:"data"`
}

func GeneratePagination(params PaginationParam) PaginationResult {
	totalPage := int(math.Ceil(float64(params.Count) / float64(params.Limit)))

	var (
		nextPage int
		PrevPage int
	)

	if params.Page < totalPage {
		nextPage = params.Page + 1
	}

	if params.Page > 1 {
		PrevPage = params.Page - 1
	}

	result := PaginationResult{
		TotalPage: totalPage,
		TotalData: params.Count,
		NextPage:  &nextPage,
		PrevPage:  &PrevPage,
		Page:      params.Page,
		Limit:     params.Limit,
		Data:      params.Data,
	}

	return result
}

func GenerateSHA256(inputString string) string {
	hash := sha256.New()
	hash.Write([]byte(inputString))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}

func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}

	return nil
}

func SetEnvfromConsul(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}

	for k, v := range env {
		valOf := reflect.ValueOf(v)
		var val string

		switch valOf.Kind() { // Perbaikan pada Kind()
		case reflect.String:
			val = valOf.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = strconv.FormatInt(valOf.Int(), 10)
		case reflect.Float32, reflect.Float64:
			val = strconv.FormatFloat(valOf.Float(), 'f', -1, 64) // Menggunakan FormatFloat
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		default:
			logrus.Errorf("unsupported type for key: %s", k) // Perbaikan log
			continue
		}

		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("failed to set env: %v", err)
			return err
		}
	}

	return nil
}

func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()
	v.SetConfigType("json")

	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		logrus.Errorf("failed to add remote provider: %v", err)
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("failed to read remote config: %v", err) // Perbaikan ejaan
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}

	err = SetEnvfromConsul(v) // Memanggil fungsi yang sudah diperbaiki
	if err != nil {
		logrus.Errorf("failed to set env from consul kv: %v", err)
		return err
	}

	return nil
}
