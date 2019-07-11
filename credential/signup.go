package credential

import (
	"net/http"
	"text/template"

	"adrdn/dit/config"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

const addNewUser = "INSERT INTO user(name, username, password) VALUES (?, ?, ?)"

// User represents the user structure
type User struct {
	ID			int
	Name		string
	Username	string
	Password	string
}

var tmpl = template.Must(template.ParseGlob("forms/user/*"))

// SignUp represent the sign up page
func SignUp(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "SignUp", nil)
}

// RegisterNewUser registers the user
func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		username := r.FormValue("username")
		password := r.FormValue("password")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			//http.Error(w, "Server error, unable to create your account.", 500)
			panic(err)
		}
		password = string(hashedPassword)
		newUserData, err := db.Prepare(addNewUser)
		if err != nil {
			panic(err)
		}
		_, err = newUserData.Exec(name, username, password)
		if err != nil {
			panic(err)
		}
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Login", nil)
}