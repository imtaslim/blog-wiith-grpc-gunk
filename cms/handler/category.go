package handler

import (
	"net/http"
	"strconv"

	cgp "blog-gunk/gunk/v1/category"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Category struct{
	ID int64
	Name string
	Status bool
	Errors map[string]string
}

func (c *Category) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name,
			validation.Required.Error("The Name Field is Required"),
			validation.Length(3, 0).Error("The Name field must be greater than or equals 3"),
		),
	)
}

func (h *Handler) categories (rw http.ResponseWriter, r *http.Request) {
	data, err := h.csc.Gets(r.Context(), &cgp.GetsCategoryRequest{})
	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(rw, "list-categories.html", data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) categoriesCreate (rw http.ResponseWriter, r *http.Request) {
	vErrs := map[string]string{}
	h.createCategoryData(rw, "", vErrs)
}

func (h *Handler) categoriesStore (rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var category Category
	if err := h.decoder.Decode(&category, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := category.Validate(); err != nil {
		valError, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range valError {
				vErrs[key] =value.Error()
			}
			h.createCategoryData(rw, category.Name, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	_, err := h.csc.Create(r.Context(), &cgp.CreateCategoryRequest{
		Category: &cgp.Category{
			Name: category.Name,
		},
	})
	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/categories", http.StatusTemporaryRedirect)
}

func (h *Handler) createCategoryData (rw http.ResponseWriter, name string, errs map[string]string) {
	form := Category{
		Name: name,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "create-categories.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) categoriesEdit (rw http.ResponseWriter, r *http.Request) {
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

	cat, err := h.csc.Get(r.Context(), &cgp.GetCategoryRequest{
		ID: id,
	})

	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}
	
	rerrs := map[string]string{}
	h.editCategoryData(rw, cat.Category.GetID(), cat.Category.Name, rerrs)
}

func (h *Handler) categoriesUpdate (rw http.ResponseWriter, r *http.Request) {
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

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var category Category
	if err := h.decoder.Decode(&category, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := category.Validate(); err != nil {
		valError, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range valError {
				vErrs[key] =value.Error()
			}
			h.editCategoryData(rw, id, category.Name, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	_, err = h.csc.Update(r.Context(), &cgp.UpdateCategoryRequest{
		Category: &cgp.Category{
			ID: id,
			Name: category.Name,
		},
	})
	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/categories", http.StatusTemporaryRedirect)
}

func (h *Handler) editCategoryData (rw http.ResponseWriter, id int64, name string, errs map[string]string) {
	form := Category{
		ID:     id,
		Name:   name,
		Errors: errs,
	}
	if err := h.templates.ExecuteTemplate(rw, "edit-categories.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) categoriesDelete (rw http.ResponseWriter, r *http.Request) {
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
	
	_, err = h.csc.Delete(r.Context(), &cgp.DeleteCategoryRequest{
		ID: id,
	})
	
	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/categories", http.StatusTemporaryRedirect)
}