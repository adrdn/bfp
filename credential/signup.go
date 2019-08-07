package credential

import (
	"net/http"

	"adrdn/dit/config"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

// SignUp represent the sign up page
func SignUp(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	var role string
	var roleList []string

	res, err := db.Query(listAllRoles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for res.Next() {
		err = res.Scan(&role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		roleList = append(roleList, role)
	}
	
	// ok, user := CheckAuthentication(w, r)
	// if !ok {
	// 	return
	// }
	tmpl.ExecuteTemplate(w, "SignUp", roleList)
}

// RegisterNewUser registers the user
func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		username := r.FormValue("username")
		password := r.FormValue("password")
		role := r.FormValue("role")
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
		_, err = newUserData.Exec(name, username, password, role)
		if err != nil {
			panic(err)
		}
	}
	defer db.Close()
	http.Redirect(w, r, "/login", 301)
}

// ChangePassword revokes the change password page
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Password", nil)
}

// UpdatePassword changes the password
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()

	if r.Method == "POST" {
		pass := r.FormValue("password")
		res, err := db.Prepare(chanPassword)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = res.Exec(pass)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	defer db.Close()
	http.Redirect(w, r, "/home", http.StatusAccepted)
}