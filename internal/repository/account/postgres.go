package account

import (
	"context"
	"fmt"

	modelRepository "github.com/AndryHardiyanto/dealltest/internal/model/repository"
	"github.com/AndryHardiyanto/dealltest/lib/errors"
	libPostgres "github.com/AndryHardiyanto/dealltest/lib/postgres"
)

type postgres struct {
	repo libPostgres.Postgres
}

//NewPostgres ..
func NewPostgres(repo libPostgres.Postgres) Repository {
	return &postgres{
		repo: repo,
	}
}

type Repository interface {
	Select(ctx context.Context, email string) (*modelRepository.Account, error)
	SelectById(ctx context.Context, accountID int64) (*modelRepository.Account, error)
	Insert(ctx context.Context, req *modelRepository.Account) error
	Delete(ctx context.Context, accountID int64) error
	Update(ctx context.Context, req *modelRepository.UpdateRequest) error
}

func (p *postgres) Update(ctx context.Context, req *modelRepository.UpdateRequest) error {
	var kv []interface{}

	querySet := ""
	kv = append(kv, "accID")
	kv = append(kv, req.AccountID)

	if req.Email != "" {
		if querySet == "" {
			querySet = "SET ddacc_email = :email"
		}
		kv = append(kv, "email")
		kv = append(kv, req.Email)
	}
	if req.Name != "" {
		if querySet == "" {
			querySet = "SET ddacc_name = :name"
		} else {
			querySet += ",ddacc_name = :name "
		}
		kv = append(kv, "name")
		kv = append(kv, req.Name)
	}
	if req.Password != "" {
		if querySet == "" {
			querySet = "SET ddacc_password = :pass"
		} else {
			querySet += ",ddacc_password = :pass "
		}
		kv = append(kv, "pass")
		kv = append(kv, req.Password)
	}
	if req.Role != "" {
		if querySet == "" {
			querySet = "SET ddacc_role = :role"
		} else {
			querySet += ",ddacc_role = :role "
		}
		kv = append(kv, "role")
		kv = append(kv, req.Role)
	}
	query := fmt.Sprintf(`UPDATE dd_account %s
	WHERE ddacc_id = :accID;`, querySet)

	_, err := p.repo.Update(query, kv...).Exec(ctx)
	if err != nil {
		return errors.NewWrapError(err, "error Update").SetType(errors.TypeInternalServerError)
	}

	return nil
}
func (p *postgres) Delete(ctx context.Context, accountID int64) error {
	query := `DELETE FROM dd_account WHERE ddacc_id = :accID`

	_, err := p.repo.Delete(query, "accID", accountID).Exec(ctx)
	if err != nil {
		return errors.NewWrapError(err, "error Delete").SetType(errors.TypeInternalServerError)
	}

	return nil
}
func (p *postgres) Insert(ctx context.Context, req *modelRepository.Account) error {
	query := `
	INSERT INTO dd_account 
	(ddacc_name, ddacc_email, ddacc_password, ddacc_role)
	VALUES(:name,:email,:pass,:role) 
	RETURNING ddacc_id;
	`

	_, err := p.repo.Insert(query, "name", req.Name, "email", req.Email, "pass", req.Password, "role", req.Role).Exec(ctx)
	if err != nil {
		return errors.NewWrapError(err, "error Insert").SetType(errors.TypeInternalServerError)
	}

	return nil
}

func (p *postgres) SelectById(ctx context.Context, accountID int64) (*modelRepository.Account, error) {
	query := "select ddacc_id as id, ddacc_email as email,ddacc_name as  name, ddacc_password as password,ddacc_role as role from dd_account where ddacc_id = :accID"
	dest := &modelRepository.Account{}

	found, err := p.repo.Select(query, dest, "accID", accountID).One(ctx)
	if err != nil {
		return nil, errors.NewWrapError(err, "error select").SetType(errors.TypeInternalServerError)
	}

	if !found {
		return nil, nil
	}

	return dest, nil
}

func (p *postgres) Select(ctx context.Context, email string) (*modelRepository.Account, error) {
	query := "select ddacc_id as id, ddacc_password as password,ddacc_role as role from dd_account where ddacc_email = :email"
	dest := &modelRepository.Account{}

	found, err := p.repo.Select(query, dest, "email", email).One(ctx)
	if err != nil {
		return nil, errors.NewWrapError(err, "error select").SetType(errors.TypeInternalServerError)
	}

	if !found {
		return nil, nil
	}

	return dest, nil
}
