package main

import (
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/tabvn/orm"
)

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      *string   `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PostID    string `json:"post_id" orm:"index"`
	Post      *Post  `json:"post" orm:"fk:PostID"`
	CreatedAt time.Time
}

func main() {
	// Open Connection
	db, err := orm.Open("postgres", "host=localhost port=5432 user=postgres dbname=orm password=root sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	// Register models, auto create table, primary key, index
	db.Models(&User{}, &Post{})
}
