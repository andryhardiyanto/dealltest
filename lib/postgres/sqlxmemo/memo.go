package sqlxmemo

import (
	"context"
	"strings"

	"github.com/AndryHardiyanto/dealltest/lib/log"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jmoiron/sqlx"
)

// SQLXMemoization is the API for github.com/jmoiron/sqlx prepared statement memoization strategy.
type SQLXMemoization interface {
	// Prepare checks cache for memoized statement by its plain query and returns it.
	// If no cache found then it will create new prepared statement, save into cache, and then returns it.
	Prepare(ctx context.Context, db *sqlx.DB, query string) (MemoizedStatement, error)
	// Prepare checks cache for memoized statement by its plain query and returns it.
	// If no cache found then it will create new prepared statement, save into cache, and then returns it.
	PrepareNamed(ctx context.Context, db *sqlx.DB, query string) (MemoizedStatement, error)
	// Purge closes all cached statements.
	Purge()
}

// New create a SQLXMemoization instance.
func New(size int) SQLXMemoization {
	l, _ := lru.NewWithEvict(size, stmtEvictionStrategy)
	return &sqlxMemoizationImpl{
		storage: l,
	}
}

func stmtEvictionStrategy(key interface{}, value interface{}) {
	memo, ok := value.(MemoizedStatement)
	if !ok {
		return
	}
	if memo.Stmt != nil {
		if err := memo.Stmt.Close(); err != nil {
			log.Warn().Err(err).Msg(err.Error())
		}
	}
	if memo.NamedStmt != nil {
		if err := memo.NamedStmt.Close(); err != nil {
			log.Warn().Err(err).Msg(err.Error())
		}
	}
}

// MemoizedStatement is a wrapper that holds the memoized statement.
type MemoizedStatement struct {
	Stmt      *sqlx.Stmt
	NamedStmt *sqlx.NamedStmt
	Query     string
}

// sqlxMemoizationImpl is the implementation for SQLXMemoization.
type sqlxMemoizationImpl struct {
	storage *lru.Cache
}

var _ SQLXMemoization = (*sqlxMemoizationImpl)(nil)

// Prepare checks cache for memoized statement by its plain query and returns it.
// If no cache found then it will create new prepared statement, save into cache, and then returns it.
func (impl *sqlxMemoizationImpl) Prepare(ctx context.Context, db *sqlx.DB, query string) (MemoizedStatement, error) {
	var empty MemoizedStatement
	key := FromString(query)
	memo, found := impl.storage.Get(key)
	existing, ok := memo.(MemoizedStatement)
	if found && ok {
		// if query is equivalent then just yield the existing
		if strings.EqualFold(existing.Query, query) {
			return existing, nil
		}

		// close existing stmt because the query is not equal
		// so that we can safely create another stmt
		_ = existing.Stmt.Close()
	}

	// create new stmt from query
	stmt, err := db.PreparexContext(ctx, query)
	if err != nil {
		return empty, err
	}

	// memo it
	existing = MemoizedStatement{
		Stmt:  stmt,
		Query: query,
	}
	impl.storage.Add(key, existing)

	return existing, err
}

// Prepare checks cache for memoized statement by its plain query and returns it.
// If no cache found then it will create new prepared statement, save into cache, and then returns it.
func (impl *sqlxMemoizationImpl) PrepareNamed(ctx context.Context, db *sqlx.DB, query string) (MemoizedStatement, error) {
	var empty MemoizedStatement
	key := FromString(query)
	memo, found := impl.storage.Get(key)
	existing, ok := memo.(MemoizedStatement)
	if found && ok {
		// if query is equivalent then just yield the existing
		if strings.EqualFold(existing.Query, query) {
			return existing, nil
		}

		// close existing stmt because the query is not equal
		// so that we can safely create another stmt
		_ = existing.NamedStmt.Close()
	}

	// create new stmt from query
	namedStmt, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		return empty, err
	}

	// memo it
	existing = MemoizedStatement{
		NamedStmt: namedStmt,
		Query:     query,
	}
	impl.storage.Add(key, existing)

	return existing, err
}

// Purge closes all cached statements.
func (impl *sqlxMemoizationImpl) Purge() {
	impl.storage.Purge()
}
