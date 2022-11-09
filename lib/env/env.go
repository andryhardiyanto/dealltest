package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetStringOrDefault(key, def string) string {
	return getEnvOrDefault(key, def)
}
func GetInt64OrDefault(key string, def int64) int64 {
	v := getEnvOrDefault(key, fmt.Sprint(def))
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return i
}
func GetTimeDurationInHourOrDefault(key string, def time.Duration) time.Duration {
	v := getEnvOrDefault(key, fmt.Sprint(def))
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return time.Duration(i * int64(time.Hour))
}
func GetTimeDurationInSecondOrDefault(key string, def time.Duration) time.Duration {
	v := getEnvOrDefault(key, fmt.Sprint(def))
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return time.Duration(i * int64(time.Second))
}
func GetBoolOrDefault(key string, def bool) bool {
	v := getEnvOrDefault(key, fmt.Sprint(def))
	i, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return i
}
func getEnvOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
