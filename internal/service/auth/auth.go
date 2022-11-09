package auth

import (
	"context"
	"encoding/json"
	"time"

	modelAuth "github.com/AndryHardiyanto/dealltest/internal/model/auth"
	"github.com/AndryHardiyanto/dealltest/lib/errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type service struct {
	signedSecret            string
	accessTokenExpDuration  time.Duration
	refreshTokenExpDuration time.Duration
}

//go:generate mockgen -destination=mock/service.go -package=mock_auth . Service
type Service interface {
	Validate(ctx context.Context, req *modelAuth.ValidateRequest) error
	GenerateJwt(ctx context.Context, subjectID, role string) (*modelAuth.Jwt, error)
	RenewJwt(ctx context.Context, req *modelAuth.RenewJwtRequest) (*modelAuth.RenewJwtResponse, error)
	GetClaims(ctx context.Context, req *modelAuth.GetClaimsRequest) (*modelAuth.Claims, error)
}

func NewService(secretKey string, accessExpDuration, refreshTokenExpDuration time.Duration) Service {
	return &service{
		signedSecret:            secretKey,
		accessTokenExpDuration:  accessExpDuration,
		refreshTokenExpDuration: refreshTokenExpDuration,
	}
}
func (s *service) RenewJwt(ctx context.Context, req *modelAuth.RenewJwtRequest) (*modelAuth.RenewJwtResponse, error) {
	err := s.Validate(ctx, &modelAuth.ValidateRequest{
		Token: req.RefreshToken,
	})
	if err != nil {
		return nil, errors.NewWrapError(err, "error validate")
	}

	dataClaims, err := s.GetClaims(ctx, &modelAuth.GetClaimsRequest{
		Token: req.RefreshToken,
	})
	if err != nil {
		return nil, errors.NewWrapError(err, "error GetClaims")
	}

	jwt, err := s.GenerateJwt(ctx, dataClaims.Subject, dataClaims.Role)
	if err != nil {
		return nil, errors.NewWrapError(err, "error GenerateJwt")
	}

	return &modelAuth.RenewJwtResponse{
		AccessToken:  jwt.AccessToken,
		AccessExp:    jwt.AccessExp,
		RefreshToken: jwt.RefreshToken,
		RefreshExp:   jwt.RefreshExp,
	}, nil
}

func (s *service) GenerateJwt(ctx context.Context, subjectID, role string) (*modelAuth.Jwt, error) {
	now := time.Now()
	dataAccessTokenClaims := modelAuth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(s.accessTokenExpDuration).Unix(),
			Id:        uuid.NewString(),
			Subject:   subjectID,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
		Role: role,
	}
	accessClams := jwt.NewWithClaims(jwt.SigningMethodHS512, dataAccessTokenClaims)

	accessToken, err := accessClams.SignedString([]byte(s.signedSecret))
	if err != nil {
		return nil, errors.NewWrapError(err, "failed signed secrect in access token").SetType(errors.TypeInternalServerError)
	}

	dataRefreshTokenClaims := modelAuth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(s.refreshTokenExpDuration).Unix(),
			Id:        uuid.NewString(),
			Subject:   subjectID,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
		Role: role,
	}
	refreshClams := jwt.NewWithClaims(jwt.SigningMethodHS512, dataRefreshTokenClaims)

	refreshToken, err := refreshClams.SignedString([]byte(s.signedSecret))
	if err != nil {
		return nil, errors.NewWrapError(err, "failed signed secrect in access token").SetType(errors.TypeInternalServerError)
	}

	return &modelAuth.Jwt{
		AccessToken:  accessToken,
		AccessExp:    time.Unix(dataAccessTokenClaims.IssuedAt, 0).Add(s.accessTokenExpDuration).Unix(),
		RefreshToken: refreshToken,
		RefreshExp:   time.Unix(dataRefreshTokenClaims.IssuedAt, 0).Add(s.accessTokenExpDuration).Unix(),
	}, nil
}

func (s *service) Validate(ctx context.Context, req *modelAuth.ValidateRequest) error {
	parseToken, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewError("Failed Parse jwt token").SetType(errors.TypeUnauthorized)
		}
		return []byte(s.signedSecret), nil
	})
	if err != nil {
		return errors.NewWrapError(err, "error jwt.parse").SetType(errors.TypeUnauthorized)
	}

	claims, ok := parseToken.Claims.(jwt.MapClaims)
	if ok {
		ok := claims.VerifyNotBefore(time.Now().Unix(), true)
		if !ok {
			return errors.NewError("validate VerifyNotBefore").SetType(errors.TypeUnauthorized)
		}

		ok = claims.VerifyExpiresAt(time.Now().Unix(), true)
		if !ok {
			return errors.NewError("validate VerifyExpiresAt").SetType(errors.TypeUnauthorized)
		}

		return nil
	}

	return errors.NewError("error format claims of jwt").SetType(errors.TypeUnauthorized)
}
func (s *service) GetClaims(ctx context.Context, req *modelAuth.GetClaimsRequest) (*modelAuth.Claims, error) {
	parseToken, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewError("Failed Parse jwt token").SetType(errors.TypeUnauthorized)
		}
		return []byte(s.signedSecret), nil
	})
	if err != nil {
		return nil, errors.NewWrapError(err, "Failed Parse jwt token").SetType(errors.TypeUnauthorized)
	}

	mapClaims, ok := parseToken.Claims.(jwt.MapClaims)
	if !ok {
		errors.NewError("error format claims of jwt").SetType(errors.TypeUnauthorized)
	}

	var claims *modelAuth.Claims
	bytes, _ := json.Marshal(mapClaims)
	json.Unmarshal(bytes, &claims)

	return claims, nil
}
