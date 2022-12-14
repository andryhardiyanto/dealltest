package config

import (
	"time"

	libEnv "github.com/AndryHardiyanto/dealltest/lib/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   Server
	Jwt      Jwt
	Logger   Logger
	Database Database
}
type Jwt struct {
	SignedSecret       string
	AccessExpDuration  time.Duration
	RefreshExpDuration time.Duration
}

type Server struct {
	Port    string
	GinMode string
}
type Logger struct {
	Name  string
	Debug bool
}

type Database struct {
	DatabaseConnection string
}

var Cfg Config

func RegisterConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	Cfg = Config{
		Server: Server{
			Port:    libEnv.GetStringOrDefault("DEALL_SERVER_PORT", ""),
			GinMode: libEnv.GetStringOrDefault("DEALL_SERVER_GIN_MODE", ""),
		},
		Database: Database{
			DatabaseConnection: libEnv.GetStringOrDefault("DEALL_DATABASE_CONNECTION", ""),
		},
		Jwt: Jwt{
			SignedSecret:       libEnv.GetStringOrDefault("DEALL_JWT_SIGNED_SECRET", ""),
			AccessExpDuration:  libEnv.GetTimeDurationInHourOrDefault("DEALL_JWT_ACCESS_EXP_DURATION", 0),
			RefreshExpDuration: libEnv.GetTimeDurationInHourOrDefault("DEALL_JWT_REFRESH_EXP_DURATION", 0),
		},
		Logger: Logger{
			Debug: libEnv.GetBoolOrDefault("DEALL_LOGGER_DEBUD", false),
		},
	}
}
