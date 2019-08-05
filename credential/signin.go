package credential

import (
	"net/http"

	"adrdn/dit/config"
	"adrdn/dit/user"
	"adrdn/dit/role"

	"golang.org/x/crypto/bcrypt"
)

const listOneUser = "SELECT name, password, role FROM user where username = ?"
const listAllRoles = "SELECT Name FROM role"
const addNewUser = "INSERT INTO user(name, username, password, role) VALUES (?, ?, ?, ?)"

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

	var name 			string
	var hashedPassword 	string
	var roleName		string

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
			err = checkData.Scan(&name, &hashedPassword, &roleName)
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
			Name:			name,
			Username:		username,
			Role: role.Role{
							Name:	roleName,
				},
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
