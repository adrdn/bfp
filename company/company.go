package company

import (
	"net/http"
	"text/template"
	"drdn/bfp/config"
)

const echoAllCompanies 	= "SELECT * FROM company"
const echoOneCompany   	= "SELECT * FROM company WHERE id = ?"
const updateCompany	   	= "UPDATE company SET name = ?, type = ? WHERE id = ?"
const addNewCompany		= "INSERT INTO company (name, type) VALUES (?, ?)"
const deleteCompany		= "DELETE FROM company WHERE id = ?"

// Company represents the company table structure
type Company struct {
	ID 		int
	Name 	string
	Type	string
}

var tmpl = template.Must(template.ParseGlob("forms/*"))
var db 	 = config.DbConn()

// ShowAllCompanies displays all of the companies
func ShowAllCompanies(w http.ResponseWriter, r *http.Request) {
	selDB, err := db.Query(echoAllCompanies)
	if err != nil {
		panic(err)
	}

	com 	:= Company{}
	comList := []Company{}

	for selDB.Next() {
		var _id 	int
		var _name 	string
		var _type	string
		err = selDB.Scan(&_id, &_name, &_type)
		if err != nil {
			panic(err)
		}
		com.ID		= _id
		com.Name 	= _name
		com.Type 	= _type
		comList 	= append(comList, com)
	}
	tmpl.ExecuteTemplate(w, "Company", comList)
	defer db.Close()
}

// Edit edits the entity
func Edit (w http.ResponseWriter, r *http.Request) {
	rowID := r.URL.Query().Get("id")
	selDB, err := db.Query(echoOneCompany, rowID)
	if err != nil {
		panic(err)
	}
	
	com := Company{}

	for selDB.Next() {
		var _id   	int
		var _name 	string
		var _type 	string

		err = selDB.Scan(&_id, &_name, &_type)
		if err != nil {
			panic(err)
		}
		com.ID 		= _id
		com.Name 	= _name
		com.Type 	= _type
	}
	tmpl.ExecuteTemplate(w, "Edit", com)
	defer db.Close()
}

// Update updates the selected entity with given data
func Update (w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		category := r.FormValue("type")
		updatedData, err := db.Prepare(updateCompany)
		if err != nil {
			panic(err)
		}
		updatedData.Exec(name, category)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// New represents the new entity page
func New (w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Insert adds the new entity
func Insert (w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		category := r.FormValue("type")
		newData, err := db.Prepare(addNewCompany)
		if err != nil {
			panic(err)
		}
		newData.Exec(name, category)
		defer db.Close()
		http.Redirect(w, r, "/", 301)
	}
}

// Delete drops the entity
func Delete (w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	deletedData, err := db.Prepare(deleteCompany)
	if err != nil {
		panic(err)
	}
	deletedData.Exec(id)
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
