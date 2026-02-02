package env

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetStrFromEnv(key string) string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		log.Panicf("%s is not set in environment", key)
	}
	return valueStr
}

func GetStrListFromEnv(key string) []string {
	valueStr := GetStrFromEnv(key)
	values := strings.Split(valueStr, ",")
	return values
}

func GetIntFromEnv(key string) int {
	valueStr := GetStrFromEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Panicf("%s couldn't be converted to int", key)
	}
	return value
}

func ReadPort(key string) int {
	port := GetIntFromEnv(key)

	if port < 1024 || port > 65353 {
		log.Panicf("Error converting environment variable <%s> to int between 1024 and 65353", key)
	}
	return port
}

func ParseDurationEnv(key string) time.Duration {
	val := GetStrFromEnv(key)
	if val == "" {
		log.Panicf("Error converting environment variable <%s> to time duration", key)
	}
	return time.ParseDuration(val)
}
