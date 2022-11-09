package postgres

import (
	"database/sql"
	"time"

	"github.com/AndryHardiyanto/dealltest/lib/log"
	"github.com/AndryHardiyanto/dealltest/lib/postgres/sqlxmemo"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	libLog "github.com/AndryHardiyanto/dealltest/lib/log"
)

type postgres struct {
	Master        *sqlx.DB
	preparedStmts sqlxmemo.SQLXMemoization
}

func NewPostgres(connectionString string) Postgres {
	return &postgres{
		Master:        connectPostgreSQL(connectionString),
		preparedStmts: sqlxmemo.New(512),
	}
}

func connectPostgreSQL(dsn string) *sqlx.DB {
	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panic().Str("dsn", dsn).Err(err).Msg(err.Error())
	}
	db := sqlx.NewDb(sqldb, "pgx")

	db.SetConnMaxLifetime(time.Duration(60) * time.Second)

	if err = db.Ping(); err != nil {
		log.Panic().Str("dsn", dsn).Err(err).Msg(err.Error())
	}

	runMigrate(sqldb)

	return db
}

func runMigrate(db *sql.DB) {
	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		libLog.Error().Msg(err.Error())
		panic(err)
	}

	mg, err := migrate.NewWithDatabaseInstance("file://db/migration", "", driver)
	if err != nil {
		libLog.Error().Msg(err.Error())
		panic(err)
	}

	err = mg.Up()
	if err != nil && err != migrate.ErrNoChange {
		libLog.Error().Msg(err.Error())
		panic(err)
	}
}
