package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewBun(viper *viper.Viper, logger *logrus.Logger) *bun.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")
	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		username, password, host, port, database,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
	))

	sqldb.SetMaxIdleConns(idleConnection)
	sqldb.SetMaxOpenConns(maxConnection)
	sqldb.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	db := bun.NewDB(sqldb, pgdialect.New())

	bun.SetLogger(logger)

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	db.AddQueryHook(newSlowQueryHook(logger, 2*time.Second))

	if err := db.Ping(); err != nil {
		logger.Fatalf("failed to connect database: %v", err)
	}

	logger.Info("Database connected successfully")
	return db
}
