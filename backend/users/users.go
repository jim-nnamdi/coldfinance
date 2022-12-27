package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jim-nnamdi/coldfinance/backend/connection"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	CreateAccount(username string, password string, email string, location string, verified int) (bool, error)
	GetUsers() ([]User, error)
	Convertpassword(hash []byte, password string) (bool, error)
	GetUserPwdHash(email string) ([]byte, error)
	GetUserByEmailAndPassword(email string, password string) (*User, error)
	LoginUser(email string, password string) (bool, error)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	EmailAdd string `json:"email"`
	Location string `json:"location"`
	Verified int    `json:"verified"`
}

func NewUser(id int, username string, password string, emailaddress string, location string, verified int) *User {
	return &User{
		Id:       id,
		Username: username,
		Password: password,
		EmailAdd: emailaddress,
		Location: location,
		Verified: verified,
	}
}

var (
	ErrCreatingUser = "could not create user account"
	dbc             = connection.Dbconn()
)

func CreateAccount(username string, password string, email string, location string, verified int) (bool, error) {
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

func GetUsers() ([]User, error) {
	res, err := dbc.Query("select * from users")
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	suser := User{}
	auser := make([]User, 0)
	for res.Next() {
		err := res.Scan(&suser.Id, &suser.Username, &suser.Password, &suser.EmailAdd, &suser.Location, &suser.Verified)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
		auser = append(auser, suser)
	}
	return auser, nil
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	allusers, err := GetUsers()
	if err != nil {
		log.Print(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allusers)
}

func Convertpassword(hash []byte, password string) (bool, error) {
	cnvtpass := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if cnvtpass != nil {
		return false, cnvtpass
	}
	return true, nil
}

func GetUserPwdHash(email string) ([]byte, error) {
	var (
		user = User{}
		err  error
	)
	res := dbc.QueryRow("select * from users where email = ?", email)
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

func GetUserByEmailAndPassword(email string, password string) (*User, error) {
	var (
		user = User{}
		err  error
	)
	req := dbc.QueryRow("select * from users where email = ?", email)
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

func LoginUser(email string, password string) (bool, error) {
	// check and convert password
	dbhash, err := GetUserPwdHash(email)
	if err != nil {
		log.Printf("cannot fetch user password: %s", string(dbhash))
		return false, err
	}

	// check & convert hashed password
	convert_password, err := Convertpassword(dbhash, password)
	if err != nil || !convert_password {
		log.Printf("cannot convert password for user: %s", string(dbhash))
		return false, err
	}
	userdata, err := GetUserByEmailAndPassword(email, string(dbhash))
	if err != nil || userdata == nil {
		log.Printf("cannot login user using: %s", email)
		return false, err
	}
	return true, nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	log.Print("reg reached here ...")
	bcryptpwd, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 14)
	if err != nil {
		log.Print(err)
		return
	}
	verified, _ := strconv.Atoi(r.FormValue("verified"))
	reg, err := CreateAccount(r.FormValue("username"), string(bcryptpwd), r.FormValue("email"), r.FormValue("location"), verified)
	if err != nil || !reg {
		if err.Error() == "this user already exists" {
			json.NewEncoder(w).Encode(err.Error())
		}
		log.Print(err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode("account created successfully!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var (
		jwt_secret = []byte("Metroboominx")
	)
	userlogin, err := LoginUser(r.FormValue("email"), r.FormValue("password"))
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, dataEncode)
	token_string, err := token.SignedString(jwt_secret)
	if err != nil {
		log.Printf("failed to generate token: %s", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "coldfinance-user",
		Value:   token_string,
		Expires: time.Now().Add(60 * time.Minute),
	})

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token_string)
}
