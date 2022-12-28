package content

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jim-nnamdi/coldfinance/backend/connection"

	"go.uber.org/zap"
)

// type Post interface {
// 	GetPosts() ([]Posts, error)
// 	GetSinglePost(slug string) (*Posts, error)
// 	AddPost() (bool, error)
// }

type Postx struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Slug     string `json:"slug"`
	Author   string `json:"author"`
	Image    string `json:"image,omitempty"`
	Approved int    `json:"approved"`
}

var (
	conn           = connection.Dbconn()
	coldfinancelog = connection.Coldfinancelog()
)

func GetPosts() ([]Postx, error) {
	res, err := conn.Query("select * from posts where approved = 1")
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	spost := Postx{}
	apost := make([]Postx, 0)
	for res.Next() {
		err := res.Scan(&spost.Id, &spost.Title, &spost.Body, &spost.Slug, &spost.Author, &spost.Image, &spost.Approved)
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
	); err != nil {
		coldfinancelog.Debug("error fetching & scanning posts", zap.String("error", err.Error()))
		return nil, err
	}
	coldfinancelog.Debug(&postmodel)
	return &postmodel, nil
}

func AddPost(title string, body string, author string) (bool, error) {
	split_title := strings.Split(title, " ")
	genslug := strings.Join(split_title, "-")
	addpost, err := conn.Exec("insert into posts(title, body, slug, author, image, approved) values(?,?,?,?,?,?)", title, body, genslug, author, "", 0)
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
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(singlepost)
}

func AddNewPost(w http.ResponseWriter, r *http.Request) {
	newpost, err := AddPost(r.FormValue("title"), r.FormValue("body"), r.FormValue("author"))
	if err != nil || !newpost {
		coldfinancelog.Debug("cannot add new post", zap.Any("error", err))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode("new post added successfully!")
}
