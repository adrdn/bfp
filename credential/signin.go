package credential

import (
	"net/http"

	"adrdn/dit/config"

	"golang.org/x/crypto/bcrypt"
)

const listOneUser = "SELECT password FROM user where username = ?"

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
			panic(err)
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
			http.Redirect(w, r, "/home", 301)
		}
		
	}
}

// Home revokes the home page
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Home", nil)
}