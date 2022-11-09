package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AndryHardiyanto/dealltest/lib/postgres/sqlxmemo"
	"github.com/pkg/errors"
)

type Postgres interface {
	Select(query string, dest interface{}, kv ...interface{}) Select
	Insert(query string, kv ...interface{}) Exec
	Update(query string, kv ...interface{}) Exec
	Delete(query string, kv ...interface{}) Exec
	FromResult(from string) string
}

type Select interface {
	One(ctx context.Context) (found bool, err error)
	Many(ctx context.Context) (found bool, err error)
}

type Exec interface {
	Exec(ctx context.Context) (int64, error)
	ExecInTx(ctx context.Context) (res *ExecResult, err error)
	Insert(query string, kv ...interface{}) Exec
	Update(query string, kv ...interface{}) Exec
	Delete(query string, kv ...interface{}) Exec
	Wrap(exec Exec) Exec
}

type ExecResult struct {
	ids map[string]int64
}

func (e *ExecResult) TxResult(query string) int64 {
	return e.ids[query]
}

type selectQuery struct {
	p        *postgres
	query    string
	kv       []interface{}
	dest     interface{}
	prepMemo sqlxmemo.MemoizedStatement
	arg      map[string]interface{}
}

type execQuery struct {
	p        *postgres
	query    string
	kv       []interface{}
	pipeline *pipeline
}

func newExecQuery(p *postgres, query string, kv []interface{}) Exec {
	return &execQuery{
		p:     p,
		query: query,
		kv:    kv,
		pipeline: &pipeline{
			m: map[string][]interface{}{},
			k: []string{},
		},
	}
}

func (e *execQuery) Exec(ctx context.Context) (int64, error) {
	if e.pipeline.isTrans() {
		return 0, errors.New("transaction please use ExecInTx()")
	}
	arg, err := Pairs(e.kv)
	if err != nil {
		return 0, err
	}
	debugQuery(e.query, arg)

	if queryType(e.query) == qInsert {
		return insert(ctx, e.p.Master, e.query, arg, e.p.preparedStmts)
	} else if queryType(e.query) == qDelete {
		return delete(ctx, e.p.Master, e.query, arg)
	}
	return update(ctx, e.p.Master, e.query, arg)
}

func (e *execQuery) ExecInTx(ctx context.Context) (res *ExecResult, err error) {
	if !e.pipeline.isTrans() {
		return nil, errors.New("not transaction please use Exec()")
	}
	e.pipeline.addFirstPipeline(e.query, e.kv)

	tx, err := e.p.Master.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	res, err = e.pipeline.runPipeline(ctx, tx)

	return
}

func (e *execQuery) Wrap(exec Exec) Exec {
	ex, ok := exec.(*execQuery)
	if !ok || ex == nil {
		return e
	}
	ex.query = e.pipeline.uniqueQuery(ex.query)
	ex.pipeline.addFirstPipeline(ex.query, ex.kv)
	e.pipeline.appendPipeline(ex.pipeline)
	return e
}

func (e *execQuery) Insert(query string, kv ...interface{}) Exec {
	e.pipeline.addPipeline(query, kv)
	return e
}

func (e *execQuery) Update(query string, kv ...interface{}) Exec {
	e.pipeline.addPipeline(query, kv)
	return e
}

func (e *execQuery) Delete(query string, kv ...interface{}) Exec {
	e.pipeline.addPipeline(query, kv)
	return e
}

func (q *selectQuery) One(ctx context.Context) (found bool, err error) {
	err = q.prep(ctx)
	if err != nil {
		return false, err
	}
	err = q.prepMemo.NamedStmt.GetContext(ctx, q.dest, q.arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (q *selectQuery) Many(ctx context.Context) (found bool, err error) {
	err = q.prep(ctx)
	if err != nil {
		return false, err
	}
	err = q.prepMemo.NamedStmt.SelectContext(ctx, q.dest, q.arg)
	if err != nil {
		return false, err
	}
	if q.dest == nil {
		return false, nil
	}

	return true, nil
}

func (q *selectQuery) prep(ctx context.Context) error {
	var err error
	q.arg, err = Pairs(q.kv)
	if err != nil {
		return err
	}
	debugQuery(q.query, q.arg)

	q.prepMemo, err = q.p.preparedStmts.PrepareNamed(ctx, q.p.Master, q.query)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgres) Select(query string, dest interface{}, kv ...interface{}) Select {
	return &selectQuery{
		p:     p,
		query: query,
		kv:    kv,
		dest:  dest,
	}
}

func (p *postgres) Insert(query string, kv ...interface{}) Exec {
	return newExecQuery(p, query, kv)
}

func (p *postgres) Update(query string, kv ...interface{}) Exec {
	return newExecQuery(p, query, kv)
}

func (p *postgres) Delete(query string, kv ...interface{}) Exec {
	return newExecQuery(p, query, kv)
}

func (p *postgres) FromResult(from string) string {
	return fmt.Sprintf("%s%s", qResult, from)
}
