package handler

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"

	cgp "blog-gunk/gunk/v1/category"
	pgp "blog-gunk/gunk/v1/post"
)

type Handler struct {
	templates *template.Template
	decoder   *schema.Decoder
	store     *sessions.CookieStore
	csc       cgp.CategoryServiceClient
	psc       pgp.PostServiceClient
}

func New(decoder *schema.Decoder, store *sessions.CookieStore, csc cgp.CategoryServiceClient, psc pgp.PostServiceClient) *mux.Router {
	h := &Handler{
		decoder: decoder,
		store:   store,
		csc:     csc,
		psc:     psc,
	}

	h.parseTemplates()

	r := mux.NewRouter()
	r.HandleFunc("/", h.home)
	r.HandleFunc("/single_post/{id:[0-9]+}", h.single)
	r.HandleFunc("/category/create", h.categoriesCreate)
	r.HandleFunc("/categories", h.categories)
	r.HandleFunc("/category/store", h.categoriesStore)
	r.HandleFunc("/category/{id:[0-9]+}/edit", h.categoriesEdit)
	r.HandleFunc("/category/{id:[0-9]+}/update", h.categoriesUpdate)
	r.HandleFunc("/category/{id:[0-9]+}/delete", h.categoriesDelete)

	r.HandleFunc("/post/create", h.postsCreate)
	r.HandleFunc("/posts", h.posts)
	r.HandleFunc("/post/store", h.postsStore)
	r.HandleFunc("/post/{id:[0-9]+}/edit", h.postsEdit)
	r.HandleFunc("/post/{id:[0-9]+}/update", h.postsUpdate)
	r.HandleFunc("/post/{id:[0-9]+}/delete", h.postsDelete)

	r.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./"))))

	r.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := h.templates.ExecuteTemplate(rw, "404.html", nil); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return r
}

func (h *Handler) parseTemplates() {
	h.templates = template.Must(template.ParseFiles(
		"cms/asset/templates/create-categories.html",
		"cms/asset/templates/list-categories.html",
		"cms/asset/templates/edit-categories.html",
		"cms/asset/templates/create-post.html",
		"cms/asset/templates/list-post.html",
		"cms/asset/templates/edit-post.html",
		"cms/asset/templates/404.html",
		"cms/asset/templates/home.html",
		"cms/asset/templates/single-post.html",
	))
}
