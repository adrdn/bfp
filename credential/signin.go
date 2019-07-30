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
	session, err := store.Get(r, "dit")
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
	session, err := store.Get(r, "dit")
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
	session, err := store.Get(r, "dit")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := getUser(session)
	
	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access:D")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tmpl.ExecuteTemplate(w, "Home", user.Username)
}


// Logout signs the user out
func Logout(w http.ResponseWriter, r *http.Request) {
	_, err := store.Get(r, "dit")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
}