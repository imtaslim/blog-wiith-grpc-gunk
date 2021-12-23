package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	cgp "blog-gunk/gunk/v1/category"
	pgp "blog-gunk/gunk/v1/post"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Post struct {
	ID          int64
	CatID       int64
	Title       string
	Description string
	Image       string
	CatName     string
	Category    []Category
	Errors      map[string]string
}

func (p *Post) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.CatID,
			validation.Required.Error("Select atleast one category"),
		),
		validation.Field(&p.Title,
			validation.Required.Error("The Title Field is Required"),
			validation.Length(3, 0).Error("The Title field must be greater than or equals 3"),
		),
		validation.Field(&p.Description,
			validation.Required.Error("The Description Field is Required"),
			validation.Length(3, 0).Error("The Description field must be greater than or equals 3"),
		),
	)
}

func (h *Handler) posts(rw http.ResponseWriter, r *http.Request) {
	data, err := h.psc.List(r.Context(), &pgp.GetsPostRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(rw, "list-post.html", data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) postsCreate(rw http.ResponseWriter, r *http.Request) {
	vErrs := map[string]string{}
	data, err := h.csc.Gets(r.Context(), &cgp.GetsCategoryRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	cats := []Category{}
	for _, v := range data.Category {
		cats = append(cats, Category{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	h.createPostData(rw, 0, "", "", cats, vErrs)
}

const MAX_UPLOAD_SIZE = 1024 * 10024 // 1MB

func (h *Handler) postsStore(rw http.ResponseWriter, r *http.Request) {
	data, err := h.csc.Gets(r.Context(), &cgp.GetsCategoryRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	cats := []Category{}
	for _, v := range data.Category {
		cats = append(cats, Category{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Body = http.MaxBytesReader(rw, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(rw, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("Image")

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	now := strconv.Itoa(int(time.Now().UnixNano()))
	img := "upload-*" + now + filepath.Ext(fileHeader.Filename)
	tempFile, err := ioutil.TempFile("cms/asset/uploads", img)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	tempFile.Write(fileBytes)
	imgName := tempFile.Name()

	var post Post
	if err := h.decoder.Decode(&post, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := post.Validate(); err != nil {
		valError, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range valError {
				vErrs[key] = value.Error()
			}
			h.createPostData(rw, post.CatID, post.Title, post.Description, cats, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.psc.Create(r.Context(), &pgp.CreatePostRequest{
		Post: &pgp.Post{
			CatID:       post.CatID,
			Title:       post.Title,
			Description: post.Description,
			Image:       imgName,
		},
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/posts", http.StatusTemporaryRedirect)
}

func (h *Handler) createPostData(rw http.ResponseWriter, catId int64, title string, desc string, cats []Category, errs map[string]string) {
	form := Post{
		CatID:       catId,
		Title:       title,
		Description: desc,
		Category:    cats,
		Errors:      errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "create-post.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) postsDelete(rw http.ResponseWriter, r *http.Request) {

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

	if err := os.Remove(post.Post.Image); err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	_, err = h.psc.Delete(r.Context(), &pgp.DeletePostRequest{
		ID: id,
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/posts", http.StatusTemporaryRedirect)
}

func (h *Handler) postsEdit(rw http.ResponseWriter, r *http.Request) {
	data, err := h.csc.Gets(r.Context(), &cgp.GetsCategoryRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	cats := []Category{}
	for _, v := range data.Category {
		cats = append(cats, Category{
			ID:   v.ID,
			Name: v.Name,
		})
	}

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

	rerrs := map[string]string{}
	h.editPostData(rw, id, post.Post.CatID, post.Post.Title, post.Post.Description, post.Post.Image, cats, rerrs)
}

func (h *Handler) postsUpdate(rw http.ResponseWriter, r *http.Request) {
	data, err := h.csc.Gets(r.Context(), &cgp.GetsCategoryRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	cats := []Category{}
	for _, v := range data.Category {
		cats = append(cats, Category{
			ID:   v.ID,
			Name: v.Name,
		})
	}

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

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	imgName := ""
	file, fileHeader, err := r.FormFile("Image")

	if file != nil {
		r.Body = http.MaxBytesReader(rw, r.Body, MAX_UPLOAD_SIZE)
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			http.Error(rw, "The uploaded file is too big. Please choose an file that's less than 10MB in size", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		now := strconv.Itoa(int(time.Now().UnixNano()))
		img := "upload-*" + now + filepath.Ext(fileHeader.Filename)
		tempFile, err := ioutil.TempFile("cms/asset/uploads", img)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		tempFile.Write(fileBytes)
		imgName = tempFile.Name()
		if err := os.Remove(post.Post.Image); err != nil {
			http.Error(rw, "Invalid URL", http.StatusInternalServerError)
			return
		}
	} else {
		imgName = post.Post.Image
	}

	var pst Post
	if err := h.decoder.Decode(&pst, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := pst.Validate(); err != nil {
		valError, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range valError {
				vErrs[key] = value.Error()
			}
			h.editPostData(rw, id, pst.CatID, pst.Title, pst.Description, post.Post.Image, cats, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.psc.Update(r.Context(), &pgp.UpdatePostRequest{
		Post: &pgp.Post{
			ID:          id,
			CatID:       pst.CatID,
			Title:       pst.Title,
			Description: pst.Description,
			Image:       imgName,
		},
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/posts", http.StatusTemporaryRedirect)
}

func (h *Handler) editPostData(rw http.ResponseWriter, id int64, catId int64, title string, desc string, img string, cats []Category, errs map[string]string) {
	form := Post{
		ID:          id,
		CatID:       catId,
		Title:       title,
		Description: desc,
		Image:       img,
		Category:    cats,
		Errors:      errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "edit-post.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
