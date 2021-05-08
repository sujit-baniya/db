package main

import (
	"fmt"
	"github.com/sujit-baniya/db"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Message string
}

func main() {
	db.Default(db.Config{
		Driver:   "postgres",
		Host:     "127.0.0.1",
		Username: "postgres",
		Password: "postgres",
		DBName:   "casbin",
		Port:     5432,
	})
	// paginate()
	// fullText()
	fullTextWithPagination()
}

func fullText() {
	var messages []Message
	db.DB.Scopes(db.FullTextFilterScope("messages", "this is test")).Find(&messages)
	fmt.Println(messages)
}

func fullTextWithPagination() {
	var messages []Message
	paging := db.Paging{
		Limit: 1,
	}
	db.DB.Scopes(db.FullTextFilterScope("messages", "I'm goodas"), db.PaginateScope(paging)).Find(&messages)
	fmt.Println(messages)
}

func paginate() {
	var messages []Message
	var message Message
	paging := db.Paging{
		Search:   "100 sending",
		SearchBy: "messages",
		Limit:    1,
	}
	data := db.Paginate(db.DB.Model(&message).Where("user_id = ?", 1), &messages, paging)

	fmt.Println(data.Items) // result data

	db.DB.Scopes(db.PaginateScope(paging)).Where("user_id = ?", 1).Find(&messages)
	fmt.Println(messages)
}
