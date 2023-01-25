package content

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jim-nnamdi/coldfinance/backend/connection"

	"go.uber.org/zap"
)

type Post interface {
	GetPosts() ([]Postx, error)
	GetSinglePost(slug string) (*Postx, error)
	AddPost(title string, body string, author string) (bool, error)
}

type Postx struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Slug       string `json:"slug"`
	Author     string `json:"author"`
	Image      string `json:"image,omitempty"`
	Approved   int    `json:"approved"`
	Category   string `json:"category"`
	DatePosted string `json:"dateposted,omitempty"`
}

// handle image uploads to s3 bucket ? needs aws subscription.
// because of aws costs we suspend blog image uploads for now
// this is still a research Proj with no funding.

var (
	conn           = connection.Dbconn()
	coldfinancelog = connection.Coldfinancelog()
)

func GetPosts() ([]Postx, error) {
	res, err := conn.Query("select * from posts where approved = 1 order by id desc limit 9")
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	spost := Postx{}
	apost := make([]Postx, 0)
	for res.Next() {
		err := res.Scan(
			&spost.Id,
			&spost.Title,
			&spost.Body,
			&spost.Slug,
			&spost.Author,
			&spost.Image,
			&spost.Approved,
			&spost.Category,
			&spost.DatePosted,
		)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
		apost = append(apost, spost)
	}
	return apost, nil
}

func GetSinglePost(slug string) (*Postx, error) {
	var (
		postmodel = Postx{}
		err       error
	)
	post := conn.QueryRow("select * from posts where slug = ?", slug)
	if err = post.Scan(
		&postmodel.Id,
		&postmodel.Title,
		&postmodel.Body,
		&postmodel.Slug,
		&postmodel.Author,
		&postmodel.Image,
		&postmodel.Approved,
		&postmodel.Category,
		&postmodel.DatePosted,
	); err != nil {
		coldfinancelog.Debug("error fetching & scanning posts", zap.String("error", err.Error()))
		return nil, err
	}
	coldfinancelog.Debug(&postmodel)
	return &postmodel, nil
}

func AddPost(title string, body string, image string, author string, category string, dateposted string) (bool, error) {
	split_title := strings.Split(title, " ")
	genslug := strings.Join(split_title, "-")
	addpost, err := conn.Exec("insert into posts(title, body, slug, author, image, approved,category,dateposted) values(?,?,?,?,?,?,?,?)", title, body, genslug, author, image, 1, category, dateposted)
	if err != nil {
		coldfinancelog.Debug("could not create new post", zap.Any("error", err))
		return false, err
	}
	newpost, err := addpost.LastInsertId()
	if err != nil || newpost == 0 {
		coldfinancelog.Debug("error adding post", zap.Any("error", err))
		return false, err
	}
	coldfinancelog.Info("new post created!")
	return true, nil
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	coldfinancelog.Debug("hitting this point", zap.Any("point", "hitting this point"))
	allposts, err := GetPosts()
	if err != nil || allposts == nil {
		coldfinancelog.Debug("cannot fetch posts", zap.Any("error", err))
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allposts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	singlepost, err := GetSinglePost(slug)
	if err != nil || singlepost == nil {
		coldfinancelog.Debug("cannot fetch post with slug", zap.Any("error", err))
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(singlepost)
}

func GetPostByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	req, err := connection.Dbconn().Query("select * from posts where category=? order by id desc", category)
	if err != nil {
		log.Print(err.Error())
		return
	}
	posts := make([]Postx, 0)
	for req.Next() {
		post := Postx{}
		err := req.Scan(
			&post.Id,
			&post.Title,
			&post.Body,
			&post.Slug,
			&post.Author,
			&post.Image,
			&post.Approved,
			&post.Category,
			&post.DatePosted,
		)
		if err != nil {
			log.Print(err.Error())
			return
		}
		posts = append(posts, post)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func AddNewPost(w http.ResponseWriter, r *http.Request) {
	newpost, err := AddPost(r.FormValue("title"), r.FormValue("body"), r.FormValue("image"), r.FormValue("author"), r.FormValue("category"), r.FormValue("dateposted"))
	if err != nil || !newpost {
		log.Print(err.Error())
		coldfinancelog.Error("cannot add new post", zap.Any("error", err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode("new post added successfully!")
}
