package account

import (
	"context"
	"fmt"

	modelAccount "github.com/AndryHardiyanto/dealltest/internal/model/account"
	"github.com/AndryHardiyanto/dealltest/internal/model/repository"
	repoAccount "github.com/AndryHardiyanto/dealltest/internal/repository/account"
	serviceAuth "github.com/AndryHardiyanto/dealltest/internal/service/auth"
	"github.com/AndryHardiyanto/dealltest/lib/errors"
	"github.com/AndryHardiyanto/dealltest/lib/hash"
)

type service struct {
	accountRepo repoAccount.Repository
	authService serviceAuth.Service
}

type Service interface {
	Login(ctx context.Context, req *modelAccount.LoginRequest) (*modelAccount.LoginResponse, error)
	Get(ctx context.Context, accountID int64) (*modelAccount.GetResponse, error)
	Update(ctx context.Context, req *modelAccount.UpdateRequest) error
	Register(ctx context.Context, req *modelAccount.RegisterRequest) error
	Delete(ctx context.Context, req *modelAccount.DeleteRequest, accountID int64) error
}

func NewService(accountRepo repoAccount.Repository, authService serviceAuth.Service) Service {
	return &service{
		accountRepo: accountRepo,
		authService: authService,
	}
}
func (s *service) Delete(ctx context.Context, req *modelAccount.DeleteRequest, accountID int64) error {
	if req.AccountID == accountID {
		return errors.NewError("can't delete self").SetType(errors.TypeCannotDeleteSelf)
	}

	_, err := s.Get(ctx, req.AccountID)
	if err != nil {
		return errors.NewWrapError(err, "error service Get")
	}

	err = s.accountRepo.Delete(ctx, req.AccountID)
	if err != nil {
		return errors.NewWrapError(err, "error service Delete")
	}

	return nil
}

func (s *service) Register(ctx context.Context, req *modelAccount.RegisterRequest) error {
	if req.Role != "admin" && req.Role != "user" {
		return errors.NewError("invalid role").SetType(errors.TypeInvalidRole)
	}

	if req.Role != "admin" && req.Role != "user" {
		return errors.NewError("invalid role").SetType(errors.TypeInvalidRole)
	}

	data, err := s.accountRepo.Select(ctx, req.Email)
	if err != nil {
		return errors.NewWrapError(err, "error accountRepo Insert")
	}

	if data != nil {
		return errors.NewError("user have been register").SetType(errors.TypeEmailAlreadyExists)
	}

	hash, err := hash.HashPassword(req.Password)
	if err != nil {
		return errors.NewWrapError(err, "error GenerateFromPassword").SetType(errors.TypeInternalServerError)
	}

	req.Password = hash

	err = s.accountRepo.Insert(ctx, &repository.Account{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
	})
	if err != nil {
		return errors.NewWrapError(err, "error accountRepo Insert")
	}

	return nil
}

func (s *service) Update(ctx context.Context, req *modelAccount.UpdateRequest) error {
	if req.Email == "" && req.Name == "" && req.Password == "" && req.Role == "" {
		return errors.NewError("request empty").SetType(errors.TypeUnprocessableEntity)
	}

	data, err := s.accountRepo.SelectById(ctx, req.AccountID)
	if err != nil {
		return errors.NewWrapError(err, "error accountRepo Update")
	}

	if data == nil {
		return errors.NewError("user id not found").SetType(errors.TypeUserNotFound)
	}

	data, err = s.accountRepo.Select(ctx, req.Email)
	if err != nil {
		return errors.NewWrapError(err, "error accountRepo Insert")
	}

	if data != nil {
		return errors.NewError("user have been register").SetType(errors.TypeEmailAlreadyExists)
	}
	if req.Password != "" {
		hash, err := hash.HashPassword(req.Password)
		if err != nil {
			return errors.NewWrapError(err, "error GenerateFromPassword").SetType(errors.TypeInternalServerError)
		}

		req.Password = hash
	}
	err = s.accountRepo.Update(ctx, &repository.UpdateRequest{
		AccountID: req.AccountID,
		Email:     req.Email,
		Name:      req.Name,
		Password:  req.Password,
		Role:      req.Role,
	})
	if err != nil {
		return errors.NewWrapError(err, "error accountRepo Update")
	}

	return nil
}

func (s *service) Get(ctx context.Context, accountID int64) (*modelAccount.GetResponse, error) {
	data, err := s.accountRepo.SelectById(ctx, accountID)
	if err != nil {
		return nil, errors.NewWrapError(err, "error accountRepo SelectById")
	}

	if data == nil {
		return nil, errors.NewError("user id not found").SetType(errors.TypeUserNotFound)
	}

	return &modelAccount.GetResponse{
		Email: data.Email,
		Id:    data.ID,
		Role:  data.Role,
		Name:  data.Name,
	}, nil
}

func (s *service) Login(ctx context.Context, req *modelAccount.LoginRequest) (*modelAccount.LoginResponse, error) {
	data, err := s.accountRepo.Select(ctx, req.Email)
	if err != nil {
		return nil, errors.NewWrapError(err, "error accountRepo Select")
	}
	if data == nil {
		return nil, errors.NewError("validation user not found").SetType(errors.TypeUserNotFound)
	}

	isMatch := hash.ComparePassword(data.Password, req.Password)
	if !isMatch {
		return nil, errors.NewError("validation password not match").SetType(errors.TypePasswordNotMatch)
	}

	dataJwt, err := s.authService.GenerateJwt(ctx, fmt.Sprintf("%d", data.ID), data.Role)
	if err != nil {
		return nil, errors.NewWrapError(err, "error authService GenerateJwt")
	}

	return &modelAccount.LoginResponse{
		AccessToken:  dataJwt.AccessToken,
		AccessExp:    dataJwt.AccessExp,
		RefreshToken: dataJwt.RefreshToken,
		RefreshExp:   dataJwt.RefreshExp,
	}, nil
}
