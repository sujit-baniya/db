package db

import (
	"fmt"
	"time"

	"github.com/sujit-baniya/log"
	gorm2 "github.com/sujit-baniya/log/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type Config struct {
	Driver      string `yaml:"driver" env:"DB_DRIVER"`
	Host        string `yaml:"host" env:"DB_HOST"`
	Username    string `yaml:"username" env:"DB_USER"`
	Password    string `yaml:"password" env:"DB_PASS"`
	DBName      string `yaml:"db_name" env:"DB_NAME"`
	Port        int    `yaml:"port" env:"DB_PORT"`
	MaxOpenCon int    `yaml:"connections" env:"DB_CONNECTIONS"`
	MaxIdleCon int    `yaml:"idle_connections" env:"DB_IDLE_CONNECTIONS"`
}

var DB *gorm.DB

func Default(cfg Config) error {
	db, err := New(cfg)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func New(cfg Config) (*gorm.DB, error) {
	var db *gorm.DB
	//nolint:wsl,lll
	var err error //nolint:wsl
	connectionString := ""
	gormLogger := gorm2.Logger{
		Log: &log.DefaultLogger,
	}
	newLogger := gormLogger.LogMode(logger.Info)
	switch cfg.Driver {
	case "postgres":
		connectionString = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password)
		db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
		})

	default:
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
		db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
		})
	}
	if err != nil {
		panic(err)
	}
	db.Use(
		dbresolver.Register(dbresolver.Config{}).
			SetConnMaxIdleTime(30 * time.Minute).
			SetConnMaxLifetime(1 * time.Hour).
			SetMaxIdleConns(80).
			SetMaxOpenConns(100),
	)
	return db, nil //nolint:wsl
}