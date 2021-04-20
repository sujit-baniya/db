package main

import (
	"fmt"
	"github.com/sujit-baniya/db"
	"gorm.io/gorm"
)

func main() {
	db.Default(db.Config{
		Driver:   "postgres",
		Host:     "127.0.0.1",
		Username: "postgres",
		Password: "postgres",
		DBName:   "casbin",
		Port:     5432,
	})
	type Message struct {
		gorm.Model
		Message string
	}

	var (
		dbEntity = db.DB.Where("")
		paging   = db.Paging{}
		bookList = struct {
			Items      []*Message
			Pagination *db.Pagination
		}{}
	)

	pages, err := db.Pages(&db.Param{
		DB:     dbEntity,
		Paging: &paging,
	}, &bookList.Items)
	if err != nil {
		return
	}
	bookList.Pagination = pages
	fmt.Printf("%+v\n", bookList.Items)      // result data
	fmt.Printf("%+v\n", bookList.Pagination) // result pagination
}
