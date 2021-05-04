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
		paging = db.Paging{Limit: 2}
	)
	var messages []Message
	var message Message
	data := db.Paginate(db.DB.Model(&message).Where("user_id = ?", 1), &messages, paging)

	fmt.Println(data.Items) // result data
}
