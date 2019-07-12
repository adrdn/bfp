package role

import (
	"net/http"
	"text/template"

	"adrdn/dit/config"
)

const echoAllRoles = "SELECT  ID, name FROM role"
const deleteUser = "DELETE FROM role WHERE id = ?" 

// Role represents the role structure
type Role struct {
	ID int
	Name string
}

var tmpl = template.Must(template.ParseGlob("forms/admin/role/*"))

// ShowAllRoles displays all of the roles
func ShowAllRoles (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	rList, err := db.Query("SELECT * from role")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	role := Role{}
	roleList := []Role{}

	for rList.Next() {
		err = rList.Scan(role.ID, role.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		roleList = append(roleList, role)
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Echo", roleList)
}

// // Edit edits the entity
// func Edit (w http.ResponseWriter, r *http.Request) {
// 	db := config.DbConn()
// 	rowID := r.URL.Query().Get("id")
// 	selDB, err := db.Query("SELECT * FROM courses WHERE id=?", rowID)
// 	config.CheckErr(err)

// 	course := Course{}

// 	for selDB.Next() {
// 		var _id int
// 		var _name string
// 		var _day string
// 		var _time string
	
// 		err = selDB.Scan(&_id, &_name, &_day, &_time)
// 		config.CheckErr(err)
// 		course.ID = _id
// 		course.Name = _name
// 		course.Day = _day
// 		course.Time = _time
// 	}

// 	tmpl.ExecuteTemplate(w, "Edit", course)
// 	defer db.Close()
// }

// // Update updates the entity with given data
// func Update (w http.ResponseWriter, r *http.Request) {
// 	db := config.DbConn()
// 	if r.Method == "POST" {
// 		id := r.FormValue("id")
// 		name := r.FormValue("name")
// 		day := r.FormValue("day")
// 		time := r.FormValue("time")
// 		updForm, err := db.Prepare("UPDATE courses SET name=?, day=?, time=? WHERE id=?")
// 		config.CheckErr(err)
// 		updForm.Exec(name, day, time, id)
// 	}
// 	defer db.Close()
// 	http.Redirect(w, r, "/course/", 301)
// }

// // New represents the new entity page
// func New (w http.ResponseWriter, r *http.Request) {
// 	tmpl.ExecuteTemplate(w, "New", nil)
// }

// // Insert adds the new entity
// func Insert (w http.ResponseWriter, r *http.Request) {
// 	db := config.DbConn()
// 	if r.Method == "POST" {
// 		Name := r.FormValue("name")
// 		Day := r.FormValue("day")
// 		Time := r.FormValue("time")
// 		insForm, err := db.Prepare("INSERT INTO courses (name, day, time) VALUES (?, ?, ?)")
// 		config.CheckErr(err)
// 		insForm.Exec(Name, Day, Time)
// 	}
// 	defer db.Close()
// 	http.Redirect(w, r, "/course/", 301)
// }

// // Delete deletes the entity
// func Delete (w http.ResponseWriter, r *http.Request) {
// 	db := config.DbConn()
// 	id := r.URL.Query().Get("id")
// 	delForm, err := db.Prepare("Delete from courses WHERE id = ?")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	delForm.Exec(id)
// 	defer db.Close()
// 	http.Redirect(w, r, "/course/", 301)
// }