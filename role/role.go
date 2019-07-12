package role

import (
	"net/http"
	"text/template"

	"adrdn/dit/config"
)

const echoAllRoles	= "SELECT  ID, name FROM role"
const echoOneRole	= "SELECT * FROM role WHERE id = ?"
const updateRole 	= "UPDATE role SET name = ? WHERE id = ?"
const newRole		= "INSERT INTO role (name) VALUES (?)"
const deleteUser 	= "DELETE FROM role WHERE id = ?" 

// Role represents the role structure
type Role struct {
	ID int
	Name string
}

var tmpl = template.Must(template.ParseGlob("forms/admin/role/*"))

// ShowAllRoles displays all of the roles
func ShowAllRoles (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	rList, err := db.Query(echoAllRoles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	role := Role{}
	roleList := []Role{}

	for rList.Next() {
		err = rList.Scan(&role.ID, &role.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		roleList = append(roleList, role)
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Echo", roleList)
}

// Edit revoke the edit page
func Edit (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	rowID := r.URL.Query().Get("id")
	rList, err := db.Query(echoOneRole, rowID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	role := Role{}

	for rList.Next() {
		err = rList.Scan(&role.ID, &role.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Edit", role)
}

// Update revises the role with the given data
func Update (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		updForm, err := db.Prepare(updateRole)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		updForm.Exec(name, id)
	}
	defer db.Close()
	http.Redirect(w, r, "/admin/role", 301)
}

// New represents the new role page
func New (w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Insert adds the new role
func Insert (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		insForm, err := db.Prepare(newRole)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		insForm.Exec(name)
	}
	defer db.Close()
	http.Redirect(w, r, "/admin/role", 301)
}

// Delete hard-deletes the role from the database
func Delete (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	id := r.URL.Query().Get("id")
	delForm, err := db.Prepare(deleteUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	delForm.Exec(id)

	defer db.Close()
	http.Redirect(w, r, "/admin/role", 301)
}