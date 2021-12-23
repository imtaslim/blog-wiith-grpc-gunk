package handler

import (
	pgp "blog-gunk/gunk/v1/post"
	cgp "blog-gunk/gunk/v1/category"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Pagination struct{
	URL string
	PageNo int
}

type ListPost struct{
	Categories []Category
	Post []Post
	SlidePost []Post
	Search string
	Pagination []Pagination
	TotalPage int
	CurrentPage int
	PrePageURL string
	NextPageURL string
}

func (h *Handler) home(rw http.ResponseWriter, r *http.Request) {
	cat, err := h.csc.Gets(r.Context(), &cgp.GetsCategoryRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	cats := []Category{}
	for _, v := range cat.Category {
		cats = append(cats, Category{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	var p int = 1
	var errr error
	var nextPageURL string
	var prePageURL string

	page := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")
	if page != "" {
		p, errr = strconv.Atoi(page)
	}

	pageQuery := fmt.Sprintf("&search=%s", search)

	if errr != nil {
		http.Error(rw, errr.Error(), http.StatusInternalServerError)
		return
	}
	offset := 0
	limit := 5

	if p > 0 {
		offset = limit * p - limit
	}
	
	total := 0

	data, err := h.psc.PaginateSearch(r.Context(), &pgp.GetsPSRequest{
		Offset: int64(offset),
		Limit:  int64(limit),
		Search: search,
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	total = int(data.Total)

	totalPage := int(math.Ceil(float64(total)/float64(limit)))
	pagination := make([]Pagination,totalPage)
	
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("cms/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("error loading configuration: %v", err)
	}
	host, port := config.GetString("server.host"),config.GetString("server.port")
	for i:=0; i<totalPage; i++ {
		pagination[i] = Pagination{
			URL: fmt.Sprintf("http://%s:%s/?page=%d%s", host, port, i +1, pageQuery),
			PageNo: i +1,
		}
		if i + 1 == p {
			if i != 0 {
				prePageURL = fmt.Sprintf("http://%s:%s/?page=%d%s", host, port, i, pageQuery)
			}
			if i + 1 != totalPage {
				nextPageURL = fmt.Sprintf("http://%s:%s/?page=%d%s", host, port, i + 2, pageQuery)
			}
		}
	}

	posts := []Post{}
	for _, v := range data.Post {
		posts = append(posts, Post{
			ID:          v.ID,
			CatID:       v.CatID,
			Title:       v.Title,
			Description: v.Description,
			Image:       v.Image,
			CatName:     v.CatName,
		})
	}

	
	Ldata, err := h.psc.List(r.Context(), &pgp.GetsPostRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	slideposts := []Post{}
	for _, v := range Ldata.Post {
		slideposts = append(slideposts, Post{
			ID:          v.ID,
			CatID:       v.CatID,
			Title:       v.Title,
			Description: v.Description,
			Image:       v.Image,
			CatName:     v.CatName,
		})
	}

	list := ListPost{
		Categories:  cats,
		Post:        posts,
		SlidePost:   slideposts,
		Search:      search,
		Pagination:  pagination,
		TotalPage:   totalPage,
		CurrentPage: p,
		PrePageURL:  prePageURL,
		NextPageURL: nextPageURL,
	}

	if err := h.templates.ExecuteTemplate(rw, "home.html", list); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) single(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(Id, 10, 64)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := h.psc.Get(r.Context(), &pgp.GetPostRequest{
		ID: id,
	})

	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(rw, "single-post.html", post); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// func (h *Handler) home(rw http.ResponseWriter, r *http.Request) {
// 	data, err := h.psc.List(r.Context(), &pgp.GetsPostRequest{})
// 	if err != nil {
// 		http.Error(rw, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if err := h.templates.ExecuteTemplate(rw, "home.html", data); err != nil {
// 		http.Error(rw, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }