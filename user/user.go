package user

import (
	"net/http"
	"text/template"

	"adrdn/dit/config"

	_ "github.com/go-sql-driver/mysql"
)

const echoAllUsers = "SELECT  ID, name, username FROM user"
const updateUser = "UPDATE user SET name = ?, username = ?, password = ? WHERE id = ?"
const deleteUser = "DELETE FROM user WHERE id = ?"

// User represents the user structure
type User struct {
	ID       int
	Name     string
	Username string
	Password string
}

var tmpl = template.Must(template.ParseGlob("forms/admin/user/*"))

// DisplayAllUsers shows the list of users
func DisplayAllUsers(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	allUsers, err := db.Query(echoAllUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	u := User{}
	uList := []User{}

	for allUsers.Next() {
		err = allUsers.Scan(&u.ID, &u.Name, &u.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		uList = append(uList, u)
	}
	tmpl.ExecuteTemplate(w, "Echo", uList)
	defer db.Close()
}

// Edit enables the user to revise her/his own data
// func Edit (w http.ResponseWriter, r *http.Request) {
// 	db := config.DbConn()
// 	rowID := r.URL.Query().Get("id")
// 	selectedUser, err := db.Query(updateUser, rowID)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for checkData.Next() {
// 		err = checkData.Scan(&_username, &hashedPassword)
// 		if err != nil || _username == "" || hashedPassword == "" {
// 			panic(err)
// 			//http.Error(w, "Invalid username", 500)
// 		}
// 		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
// 			panic(err)
// 			//http.Error(w, "Invalid Password", 500)
// 		}
// 	}
// 	com := Company{}

// 	for selDB.Next() {
// 		var _id   	int
// 		var _name 	string
// 		var _type 	string

// 		err = selDB.Scan(&_id, &_name, &_type)
// 		if err != nil {
// 			panic(err)
// 		}
// 		com.ID 		= _id
// 		com.Name 	= _name
// 		com.Type 	= _type
// 	}
// 	tmpl.ExecuteTemplate(w, "Edit", com)
// 	defer db.Close()
// }

// // Update updates the selected entity with given data
// func Update (w http.ResponseWriter, r *http.Request) {
// 	db := config.DbConn()
// 	if r.Method == "POST" {
// 		id := r.FormValue("id")
// 		name := r.FormValue("name")
// 		category := r.FormValue("type")
// 		updatedData, err := db.Prepare(updateCompany)
// 		if err != nil {
// 			panic(err)
// 		}
// 		updatedData.Exec(name, category, id)
// 	}
// 	defer db.Close()
// 	http.Redirect(w, r, "/company", 301)
// }

// // New represents the new entity page
// func New (w http.ResponseWriter, r *http.Request) {
// 	tmpl.ExecuteTemplate(w, "New", nil)
// }

// // Insert adds the new entity
// func Insert (w http.ResponseWriter, r *http.Request) {
// 	db := config.DbConn()
// 	if r.Method == "POST" {
// 		name := r.FormValue("name")
// 		category := r.FormValue("type")
// 		newData, err := db.Prepare(addNewCompany)
// 		if err != nil {
// 			panic(err)
// 		}
// 		newData.Exec(name, category)
// 		defer db.Close()
// 		http.Redirect(w, r, "/company", 301)
// 	}
// }

// DeleteUser drops the user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	id := r.URL.Query().Get("id")
	deletedUser, err := db.Prepare(deleteUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	deletedUser.Exec(id)
	defer db.Close()
	http.Redirect(w, r, "/admin/users", 301)
}
