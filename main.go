package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jim-nnamdi/coldfinance/backend/users"
)

func main() {
	log.Print("server running on 9900 ...")
	http.HandleFunc("/users", users.GetAllUsers)
	http.HandleFunc("/register", users.Register)
	http.HandleFunc("/login", users.Login)
	http.ListenAndServe(":9900", nil)
}
