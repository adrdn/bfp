package credential

import (
	"net/http"

	"adrdn/dit/config"
	"adrdn/dit/user"

	"golang.org/x/crypto/bcrypt"
)

const listOneUser = "SELECT password FROM user where username = ?"

// Login revokes the login page
func Login(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "dit")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	tmpl.ExecuteTemplate(w, "Login", user)
}

// Authentication decides if the user can login or not
func Authentication(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	var hashedPassword string
	session, err := Store.Get(r, "dit")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
			session.AddFlash("The username or password is incorrect")
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		
		defer db.Close()
		u := user.User{
			Username:		username,
			Authenticated:	true,
		}
		session.Values["user"] = u
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

// Home revokes the home page
func Home(w http.ResponseWriter, r *http.Request) {
	ok, user := CheckAuthentication(w, r)
	if !ok {
		return
	}
	tmpl.ExecuteTemplate(w, "Home", user.Username)
}


// Logout signs the user out
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "dit")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["user"] = user.User{}
		session.Options.MaxAge = -1

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
}