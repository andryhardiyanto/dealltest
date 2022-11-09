package app

import (
	"github.com/AndryHardiyanto/dealltest/internal/service/account"
	"github.com/AndryHardiyanto/dealltest/internal/service/auth"
)

type App struct {
	Services *Services
}

type Services struct {
	AuthService    auth.Service
	AccountService account.Service
}
