package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/hbourgeot/portfolio/internal/store-crud"
	"github.com/hbourgeot/portfolio/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type templateData struct {
	Products []*crud.Products
	Form     any
	Flash    string
}

type insertForm struct {
	Code                int    `form:"code"`
	Name                string `form:"name"`
	Category            string `form:"category"`
	Price               int    `form:"price"`
	ImageUrl            string `form:"image_url"`
	validator.Validator `form:"-"`
}

type modifyForm struct {
	Code                int    `form:"code"`
	Modify              string `form:"modify"`
	NewVal              any    `form:"new_val"`
	validator.Validator `form:"-"`
}

type deleteForm struct {
	Code                int `form:"code"`
	validator.Validator `form:"-"`
}

type loginForm struct {
	User                string `form:"username"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	var form loginForm

	err := decodePostForm(r, &form)
	if err != nil {
		return
	}

	form.CheckField(validator.NotBlank(form.User), "user", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.User, 20), "username", "This field must be 20 characters long")
	form.CheckField(validator.Matches(form.User, validator.EmailRX), "username", "Please, enter a valid email")

	if !form.Valid() {
		data := &templateData{}
		data.Form = form
		temp := template.Must(template.ParseFiles("./ui/store-login.gohtml"))
		err = temp.Execute(w, data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	err = crud.ConfirmLogin(form.User, form.Password)
	if err != nil {
		http.Redirect(w, r, "/invalid-request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/store/admin/yes", http.StatusSeeOther)
}

func AdminCRUD(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	if logged := params.ByName("logged"); logged != "yes" {
		http.Redirect(w, r, "/not-found", http.StatusNotFound)
		return
	}

	products, err := crud.GetAllProducts()
	if err != nil {
		log.Fatal(err)
		return
	}

	data := &templateData{
		Products: products,
	}

	temp := template.Must(template.ParseFiles("./ui/store-panel.gohtml"))
	err = temp.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func AdminCreate(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	if logged := params.ByName("logged"); logged != "yes" {
		http.Redirect(w, r, "/not-found", http.StatusNotFound)
		return
	}

	var form insertForm

	err := decodePostForm(r, &form)
	if err != nil {
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Category), "category", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.ImageUrl), "image_url", "This field cannot be blank")

	form.CheckField(validator.MinChars(form.Name, 5), "name", "This field must be at least 5 characters long")
	form.CheckField(validator.MinChars(form.Category, 3), "category", "This field must be at least 5 characters long")
	form.CheckField(validator.MinChars(form.ImageUrl, 15), "image_url", "This field must be at least 15 characters long")

	form.CheckField(validator.MaxChars(form.Name, 250), "name", "This field cannot be more than 250 characters long")
	form.CheckField(validator.MaxChars(form.Category, 60), "category", "This field cannot be more than 60 characters long")
	form.CheckField(validator.MaxChars(form.ImageUrl, 300), "image_url", "This field cannot be more than 300 characters long")

	form.CheckField(validator.NotNegative(form.Code), "code", "The code must be positive")
	form.CheckField(validator.NotNegative(form.Price), "price", "The price must be positive")

	form.CheckField(validator.NotZero(form.Code), "code", "The code must be different than 0")
	form.CheckField(validator.NotZero(form.Price), "price", "The product cannot be free")

	products, err := crud.GetAllProducts()
	if err != nil {
		log.Fatal(err)
		return
	}

	if !form.Valid() {
		data := &templateData{}
		data.Form = form
		data.Products = products

		temp := template.Must(template.ParseFiles("./ui/store-panel.gohtml"))
		err = temp.Execute(w, data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	err = crud.InsertProducts(form.Code, form.Name, form.Category, form.ImageUrl, form.Price)
	if err != nil {
		data := &templateData{}

		form.AddFieldError("insert", err.Error())
		data.Form = form
		data.Products = products

		temp := template.Must(template.ParseFiles("./ui/store-panel.gohtml"))
		err = temp.Execute(w, data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	http.Redirect(w, r, "/store/admin/yes", http.StatusSeeOther)
}

func GetProductforCart(w http.ResponseWriter, r *http.Request) {
	var code int
	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		fmt.Println(err)
		return
	}

	product, err := crud.GetProductByCode(code)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(product)
}

func AdminUpdate(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	if logged := params.ByName("logged"); logged != "yes" {
		http.Redirect(w, r, "/not-found", http.StatusNotFound)
		return
	}

	var form modifyForm

	err := decodePostForm(r, &form)
	if err != nil {
		return
	}

	form.CheckField(validator.NotBlank(form.Modify), "modify", "This field cannot be blank")
	form.CheckField(validator.NotBlank(fmt.Sprintf("%v", form.NewVal)), "new_val", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Modify, 10), "column_modify", "This field cannot be more than 10 characters long")

	form.CheckField(validator.NotNegative(form.Code), "modCode", "The code must be positive")
	form.CheckField(validator.NotZero(form.Code), "modCode", "The code must be different than 0")

	products, err := crud.GetAllProducts()
	if err != nil {
		log.Fatal(err)
		return
	}

	if !form.Valid() {
		data := &templateData{}
		data.Form = form
		data.Products = products

		temp := template.Must(template.ParseFiles("./ui/store-panel.gohtml"))
		err = temp.Execute(w, data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	err = crud.UpdateProducts(form.Modify, form.NewVal, form.Code)
	if err != nil {
		data := &templateData{}
		form.AddFieldError("modify", err.Error())
		data.Form = form
		data.Products = products

		temp := template.Must(template.ParseFiles("./ui/store-panel.gohtml"))
		err = temp.Execute(w, data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	http.Redirect(w, r, "/store/admin/yes", http.StatusSeeOther)
}

func AdminDelete(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	if logged := params.ByName("logged"); logged != "yes" {
		http.Redirect(w, r, "/not-found", http.StatusNotFound)
		return
	}

	var form deleteForm

	err := decodePostForm(r, &form)
	if err != nil {
		return
	}

	form.CheckField(validator.NotNegative(form.Code), "delCode", "The code must be positive")
	form.CheckField(validator.NotZero(form.Code), "delCode", "The code must be different than 0")

	products, err := crud.GetAllProducts()
	if err != nil {
		log.Fatal(err)
		return
	}

	if !form.Valid() {
		data := &templateData{}
		data.Form = form
		data.Products = products

		temp := template.Must(template.ParseFiles("./ui/store-panel.gohtml"))
		err = temp.Execute(w, data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	err = crud.DeleteProducts(form.Code)
	if err != nil {
		data := &templateData{}

		form.AddFieldError("delete", err.Error())
		data.Form = form

		products, err := crud.GetAllProducts()
		if err != nil {
			log.Fatal(err)
			return
		}
		data.Products = products

		temp := template.Must(template.ParseFiles("./ui/store-panel.gohtml"))
		err = temp.Execute(w, data)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	http.Redirect(w, r, "/store/admin/yes", http.StatusSeeOther)
}

func ShowProducts(w http.ResponseWriter, r *http.Request) {
	products, err := crud.GetAllProducts()
	if err != nil {
		log.Fatal(err)
		return
	}

	data := &templateData{
		Products: products,
	}

	temp := template.Must(template.ParseFiles("./ui/products.gohtml"))
	err = temp.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func NewOrder(w http.ResponseWriter, r *http.Request) {
	var order crud.JSONorder
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Println(err)
		return
	}
	var cart []int
	for i := range order.Cart {
		product, err := strconv.Atoi(order.Cart[i])
		if err != nil {
			fmt.Println(err)
			return
		}

		cart = append(cart, product)
	}
	err = crud.CreateClient(order.UserName, order.Name, order.Country)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		fmt.Println(err)
	}

	err = crud.GenerateOrder(order.UserName, cart)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	formDecoder := form.NewDecoder()

	err = formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}
