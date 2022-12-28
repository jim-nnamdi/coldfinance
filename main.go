package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jim-nnamdi/coldfinance/backend/content"
	"github.com/jim-nnamdi/coldfinance/backend/users"
)

func main() {
	log.Print("server running on 9900 ...")
	r := http.NewServeMux()
	r.HandleFunc("/users", users.GetAllUsers)
	r.HandleFunc("/register", users.Register)
	r.HandleFunc("/login", users.Login)

	r.HandleFunc("/posts", content.GetAllPosts)
	r.HandleFunc("/post", content.GetPost)
	r.HandleFunc("/add/post", content.AddNewPost)
	err := http.ListenAndServe(":9900", r)
	if err != nil {
		log.Fatal(err)
	}
}
