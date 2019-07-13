package credential

import (
	"net/http"

	"adrdn/dit/config"

	"golang.org/x/crypto/bcrypt"
)

const listOneUser = "SELECT password FROM user where username = ?"
var Authenticated bool

// Login revokes the login page
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Login", nil)
}

// Authentication decides if the user can login or not
func Authentication(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	var hashedPassword string

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		checkData, err := db.Query(listOneUser, username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		for checkData.Next() {
			err = checkData.Scan(&hashedPassword)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		}

		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			defer db.Close()
			http.Error(w, err.Error(), 500)
		} else {
			defer db.Close()
			Authenticated = true
			http.Redirect(w, r, "/home", 301)
		}
		
	}
}

// Home revokes the home page
func Home(w http.ResponseWriter, r *http.Request) {
	if Authenticated {
		tmpl.ExecuteTemplate(w, "Home", nil)
	} else {
		http.Redirect(w, r, "/login", 301)
	}
}

// Logout signs out the user
func Logout(w http.ResponseWriter, r *http.Request) {
	Authenticated = false
	http.Redirect(w, r, "/login", 301)
}