package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/AndryHardiyanto/dealltest/lib/postgres/sqlxmemo"
	"github.com/jmoiron/sqlx"
)

func insertTx(ctx context.Context, tx *sqlx.Tx, query string, arg map[string]interface{}) (int64, error) {
	insertedID := 0
	prep, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	err = prep.GetContext(ctx, &insertedID, arg)
	if err != nil {
		return 0, err
	}
	if insertedID == 0 {
		return 0, errors.New("error data insertedId 0")
	}
	return int64(insertedID), nil
}

func updateTx(ctx context.Context, tx *sqlx.Tx, query string, arg map[string]interface{}) (int64, error) {
	res, err := tx.NamedExecContext(ctx, query, arg)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, errors.New("error update data rows affected 0")
	}
	return count, nil
}

func deleteTx(ctx context.Context, tx *sqlx.Tx, query string, arg map[string]interface{}) (int64, error) {
	res, err := tx.NamedExecContext(ctx, query, arg)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

func insert(ctx context.Context, db *sqlx.DB, query string, arg map[string]interface{}, p sqlxmemo.SQLXMemoization) (int64, error) {
	insertedID := 0
	prep, err := p.PrepareNamed(ctx, db, query)
	if err != nil {
		return 0, err
	}
	err = prep.NamedStmt.GetContext(ctx, &insertedID, arg)
	if err != nil {
		return 0, err
	}
	if insertedID == 0 {
		return 0, errors.New("error data insertedId 0")
	}
	return int64(insertedID), err
}

func update(ctx context.Context, db *sqlx.DB, query string, arg map[string]interface{}) (int64, error) {
	res, err := db.NamedExecContext(ctx, query, arg)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, errors.New("error update data rows affected 0")
	}
	return count, nil
}

func delete(ctx context.Context, db *sqlx.DB, query string, arg map[string]interface{}) (int64, error) {
	res, err := db.NamedExecContext(ctx, query, arg)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

type pipeline struct {
	m map[string][]interface{}
	k []string
}

func (p *pipeline) runPipeline(ctx context.Context, tx *sqlx.Tx) (*ExecResult, error) {
	res := &ExecResult{
		ids: map[string]int64{},
	}
	for _, q := range p.k {
		var id int64
		v := p.m[q]
		arg, err := PairsHook(v, res.ids, qResult)
		if err != nil {
			return nil, err
		}
		if len(p.k) < 2000 {
			debugQuery(q, arg)
		}

		qType := queryType(q)
		if strings.EqualFold(qType, qInsert) {
			id, err = insertTx(ctx, tx, q, arg)
		} else if strings.EqualFold(qType, qDelete) {
			id, err = deleteTx(ctx, tx, q, arg)
		} else {
			id, err = updateTx(ctx, tx, q, arg)
		}
		if err != nil {
			return nil, err
		}
		res.ids[q] = id
	}
	return res, nil
}

func (p *pipeline) addPipeline(query string, kv []interface{}) {
	query = p.uniqueQuery(query)
	p.m[query] = kv
	p.k = append(p.k, query)
}

func (p *pipeline) addFirstPipeline(query string, kv []interface{}) {
	if query == "" {
		return
	}
	query = p.uniqueQuery(query)
	p.m[query] = kv
	p.k = append([]string{query}, p.k...)
}

func (p *pipeline) appendPipeline(pip *pipeline) {
	for _, q := range pip.k {
		q = p.uniqueQuery(q)
		p.m[q] = pip.m[q]
		p.k = append(p.k, q)
	}
}

func (p *pipeline) isTrans() bool {
	return len(p.k) > 0
}

//add uniqueness to duplicate query
func (p *pipeline) uniqueQuery(query string) string {
	_, ok := p.m[query]
	if ok {
		query = fmt.Sprintf("%s/*%d*/", query, len(p.k))
	}
	return query
}
