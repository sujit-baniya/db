package db

import (
	"errors"
	"fmt"
	"math"
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
	Driver     string `yaml:"driver" env:"DB_DRIVER"`
	Host       string `yaml:"host" env:"DB_HOST"`
	Username   string `yaml:"username" env:"DB_USER"`
	Password   string `yaml:"password" env:"DB_PASS"`
	DBName     string `yaml:"db_name" env:"DB_NAME"`
	Port       int    `yaml:"port" env:"DB_PORT"`
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

type Pagination struct {
	TotalRecords int64 `json:"total_records" query:"total_records" form:"total_records"`
	TotalPage    int   `json:"total_page" query:"total_page" form:"total_page"`
	Offset       int   `json:"offset" query:"offset" form:"offset"`
	Limit        int   `json:"limit" query:"limit" form:"limit"`
	Page         int   `json:"page" query:"page" form:"page"`
	PrevPage     int   `json:"prev_page" query:"prev_page" form:"prev_page"`
	NextPage     int   `json:"next_page" query:"" form:""`
}

type Paging struct {
	Page    int      `json:"page" query:"page" form:"page"`
	OrderBy []string `json:"order_by" query:"order_by" form:"order_by"`
	Limit   int      `json:"limit" query:"limit" form:"limit"`
	ShowSQL bool
}

type Param struct {
	DB     *gorm.DB
	Paging *Paging
}

// Endpoint for pagination
func Pages(p *Param, result interface{}) (paginator *Pagination, err error) {

	var (
		done     = make(chan bool, 1)
		db       = p.DB.Session(&gorm.Session{})
		defPage  = 1
		defLimit = 20
		count    int64
		offset   int
	)

	// get all counts
	go getCounts(db, result, done, &count)

	// if not defined
	if p.Paging == nil {
		p.Paging = &Paging{}
	}

	// debug sql
	if p.Paging.ShowSQL {
		db = db.Debug()
	}
	// limit
	if p.Paging.Limit == 0 {
		p.Paging.Limit = defLimit
	}
	// page
	if p.Paging.Page < 1 {
		p.Paging.Page = defPage
	} else if p.Paging.Page > 1 {
		offset = (p.Paging.Page - 1) * p.Paging.Limit
	}
	// sort
	if len(p.Paging.OrderBy) > 0 {
		for _, o := range p.Paging.OrderBy {
			db = db.Order(o)
		}
	} else {
		str := "id desc"
		p.Paging.OrderBy = append(p.Paging.OrderBy, str)
	}

	// get
	if errGet := db.Limit(p.Paging.Limit).Offset(offset).Find(result).Error; errGet != nil && !errors.Is(errGet, gorm.ErrRecordNotFound) {
		return nil, errGet
	}
	<-done

	// total pages
	total := int(math.Ceil(float64(count) / float64(p.Paging.Limit)))

	// construct pagination
	paginator = &Pagination{
		TotalRecords: count,
		Page:         p.Paging.Page,
		Offset:       offset,
		Limit:        p.Paging.Limit,
		TotalPage:    total,
		PrevPage:     p.Paging.Page,
		NextPage:     p.Paging.Page,
	}

	// prev page
	if p.Paging.Page > 1 {
		paginator.PrevPage = p.Paging.Page - 1
	}
	// next page
	if p.Paging.Page != paginator.TotalPage {
		paginator.NextPage = p.Paging.Page + 1
	}

	return paginator, nil
}

func getCounts(db *gorm.DB, anyType interface{}, done chan bool, count *int64) {
	db.Model(anyType).Count(count)
	done <- true
}

func (p Pagination) IsEmpty() bool {
	return p.TotalRecords <= 0
}
