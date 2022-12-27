package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	EmailAdd string `json:"email"`
	Location string `json:"location"`
	Verified int    `json:"verified"`
}

var (
	ErrCreatingUser = "could not register user"
)

func dbconn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:M@etroboomin50@tcp(localhost:8889)/coldfinance")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return db, nil
}

func createAccount(username string, password string, email string, location string, verified int) (bool, error) {
	dbc, err := dbconn()
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	res, err := dbc.Exec("insert into users(username, password, email, location, verified) values(?,?,?,?,?)", username, password, email, location, verified)
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	userid, err := res.LastInsertId()
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	if userid != 0 {
		return true, nil
	}
	return false, errors.New(ErrCreatingUser)
}

func getUsers() (*[]User, error) {
	dbc, err := dbconn()
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	res, err := dbc.Query("select * from users")
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	suser := User{}
	auser := make([]User, 0)
	for res.Next() {
		err := res.Scan(&suser.Id, &suser.Username, &suser.Password, &suser.EmailAdd, &suser.Location, &suser.Password)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
	}
	auser = append(auser, suser)
	return &auser, nil
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	x, err := getUsers()
	if err != nil {
		log.Print(err)
	}
	json.NewEncoder(w).Encode(x)
}

func register(w http.ResponseWriter, r *http.Request) {
	conn, err := dbconn()
	if err != nil {
		log.Print(err.Error())
		return
	}
	bcryptpwd, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 14)
	if err != nil {
		log.Print(err)
		return
	}
	verified, _ := strconv.Atoi(r.FormValue("verified"))
	reg, err := createAccount(r.FormValue("username"), string(bcryptpwd), r.FormValue("email"), r.FormValue("location"), verified)
	if err != nil {
		log.Print(err.Error())
		return
	}
	if reg {
		var val User
		fetchuser := conn.QueryRow("select * from users where email = ?", r.FormValue("email")).Scan(val)
		if fetchuser != nil {
			log.Print(err.Error())
			return
		}
		generateUser := User{val.Id, val.Username, val.Password, val.EmailAdd, val.Location, val.Verified}
		json.NewEncoder(w).Encode(generateUser)
		return
	}
	json.NewEncoder(w).Encode(ErrCreatingUser)
}

func main() {
	log.Print("server running on 9900 ...")
	http.HandleFunc("/users", getAllUsers)
	http.HandleFunc("/register", register)
	http.ListenAndServe(":9900", nil)
}
