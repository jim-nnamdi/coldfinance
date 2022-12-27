package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
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
	db, err := sql.Open("mysql", "root:M@etroboomin50@tcp(localhost:3306)/coldfinance")
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

	// check for duplicates
	dup, err := dbc.Query("select * from users where email = ?", email)
	if err != nil {
		log.Print(err.Error())
		return false, err
	}
	suser := User{}
	for dup.Next() {
		if err = dup.Scan(&suser.Id, &suser.Username, &suser.Password, &suser.EmailAdd, &suser.Location, &suser.Verified); err != nil {
			log.Print(err.Error())
			return false, err
		}
	}

	// checks for email and username duplicates
	if suser.EmailAdd == email || suser.Username == username {
		fmt.Println("user already exists")
		return false, errors.New("this user already exists")
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

func getUsers() ([]User, error) {
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
		auser = append(auser, suser)
	}
	return auser, nil
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	x, err := getUsers()
	if err != nil {
		log.Print(err)
	}
	json.NewEncoder(w).Encode(x)
}

func convertpassword(hash []byte, password string) (bool, error) {
	cnvtpass := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if cnvtpass != nil {
		return false, cnvtpass
	}
	return true, nil
}

func getUserPwdHash(email string) ([]byte, error) {
	var (
		user = User{}
		err  error
	)
	dbx, err := dbconn()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	res := dbx.QueryRow("select * from users where email = ?", email)
	if err = res.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.EmailAdd,
		&user.Location,
		&user.Verified,
	); err != nil {
		log.Print(err)
		return nil, err
	}
	return []byte(user.Password), nil
}

func getUserByEmailAndPassword(email string, password string) (*User, error) {
	var (
		user = User{}
	)
	dbx, err := dbconn()
	if err != nil {
		log.Printf("no connection: %v", dbx)
		return nil, err
	}
	req := dbx.QueryRow("select * from users where email = ?")
	if err = req.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.EmailAdd,
		&user.Location,
		&user.Verified,
	); err != nil {
		log.Printf("cannot scan rows: %s", err)
		return nil, err
	}
	return &user, nil
}

func loginUser(email string, password string) (bool, error) {
	// check and convert password
	dbhash, err := getUserPwdHash(email)
	if err != nil {
		log.Printf("cannot fetch user password: %s", string(dbhash))
		return false, err
	}

	// check & convert hashed password
	convert_password, err := convertpassword(dbhash, password)
	if err != nil || !convert_password {
		log.Printf("cannot convert password for user: %s", string(dbhash))
		return false, err
	}
	userdata, err := getUserByEmailAndPassword(email, string(dbhash))
	if err != nil || userdata == nil {
		log.Printf("cannot login user using: %s", email)
		return false, err
	}
	return true, nil
}

func register(w http.ResponseWriter, r *http.Request) {
	log.Print("reg reached here ...")
	bcryptpwd, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 14)
	if err != nil {
		log.Print(err)
		return
	}
	verified, _ := strconv.Atoi(r.FormValue("verified"))
	reg, err := createAccount(r.FormValue("username"), string(bcryptpwd), r.FormValue("email"), r.FormValue("location"), verified)
	if err != nil || !reg {
		if err.Error() == "this user already exists" {
			json.NewEncoder(w).Encode(err.Error())
		}
		log.Print(err.Error())
		return
	}
	json.NewEncoder(w).Encode("account created successfully!")
}

func login(w http.ResponseWriter, r *http.Request) {
	userlogin, err := loginUser(r.FormValue("email"), r.FormValue("password"))
	if err != nil || !userlogin {
		log.Printf("error logging in with %s and %s", r.FormValue("email"), r.FormValue("password"))
		return
	}

	// we need to auth the user using jwt
	expires := time.Now().Add(time.Hour)
	type ToEncode struct {
		Email string `json:"email"`
		jwt.StandardClaims
	}
	dataEncode := ToEncode{
		Email: r.FormValue("email"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "coldfinance",
		},
	}
	claims := jwt.MapClaims{}
	tokenstring := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenstring.SignedString(dataEncode)
	if err != nil {
		log.Printf("failed to generate token: %v", tokenstring)
		return
	}
	json.NewEncoder(w).Encode(token)
}

func main() {
	log.Print("server running on 9900 ...")
	http.HandleFunc("/users", getAllUsers)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":9900", nil)
}
