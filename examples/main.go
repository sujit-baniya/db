package main

import (
	"context"
	"fmt"
	"github.com/sujit-baniya/db"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Message string
}

func main() {
	database, err := db.New(db.Config{
		Driver:   "postgres",
		Host:     "127.0.0.1",
		Username: "postgres",
		Password: "postgres",
		DBName:   "verify",
		Port:     5432,
	})
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	messageRepo := db.NewGormRepository[Message](ctx, database)
	fmt.Println(messageRepo.RawMapSlice(ctx, "SELECT * FROm messages LIMIT 10"))
}

func fullText() {
	var messages []Message
	db.DB.Scopes(db.FullTextFilterScope("messages", "this is test")).Find(&messages)
	fmt.Println(messages)
}

type OperationView struct {
	Name string `json:"name" gorm:"column:name"`
	File string `json:"file" gorm:"column:file"`
}

func fullTextWithPagination() {
	var messages []OperationView
	paging := db.Paging{
		Limit:  1,
		Search: "c7kjhoh3ol0msf6bqs90.csv",
	}
	db.DB.Scopes(db.PaginateScope(paging)).Find(&messages)
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
