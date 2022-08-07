package migrations

import (
	"embed"

	"gopkg.in/gorp.v1"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var dialects = map[string]gorp.Dialect{
	"sqlite3":  gorp.SqliteDialect{},
	"postgres": gorp.PostgresDialect{},
	"mysql":    gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"},
}

type Config struct {
	EmbeddedFS embed.FS
	Dir        string `yaml:"directory"`
	TableName  string `yaml:"table"`
	Dialect  string `yaml:"dialect"`
}

var DefaultConfig = Config{Dir: "./storage/database/migrations"}

func NewMigration(cfg Config) {

	if cfg.Dir == "" {
		cfg.Dir = "migrations"
	}
	if cfg.Dir == "" {
		cfg.Dialect = "postgres"
	}

	if cfg.TableName != "" {
		SetTable(cfg.TableName)
	} else {
		SetTable("")
	}
	DefaultConfig = cfg
}
