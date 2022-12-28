package content

import (
	"log"

	"github.com/jim-nnamdi/coldfinance/backend/connection"
	"go.uber.org/zap"
)

type Post interface {
	GetPosts() ([]Posts, error)
	GetSinglePost(slug string) (Posts, error)
	AddPost() (bool, error)
}

type Posts struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Slug   string `json:"slug"`
	Author string `json:"author"`
	Image  string `json:"image"`
}

var (
	conn           = connection.Dbconn()
	coldfinancelog = connection.Coldfinancelog()
)

func GetPosts() ([]Posts, error) {
	var (
		err error
	)
	posts, err := conn.Query("select * from posts order by id desc limit 15")
	if err != nil {
		log.Print("error fetching posts")
		coldfinancelog.Infof("error fetching posts: %v", posts)
		return nil, err
	}
	spost := Posts{}
	allPosts := make([]Posts, 0)
	for posts.Next() {
		if err := posts.Scan(&spost.Id, &spost.Title, &spost.Body, &spost.Slug, &spost.Author, &spost.Image); err != nil {
			coldfinancelog.Debug("cannot scan rows for posts", zap.String("error", err.Error()))
			return nil, err
		}
		allPosts = append(allPosts, spost)
	}
	coldfinancelog.Debug("posts returned", zap.Any("posts", allPosts))
	return allPosts, nil
}

func GetSinglePost(slug string) (*Posts, error) {
	var (
		postmodel = Posts{}
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
	); err != nil {
		coldfinancelog.Debug("error fetching & scanning posts", zap.String("error", err.Error()))
		return nil, err
	}
	return &postmodel, nil
}

func AddPost() (bool, error)
