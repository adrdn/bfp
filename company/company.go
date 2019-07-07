package company

import (
	"net/http"
	"text/template"
	"drdn/bfp/config"
)

const echoAllCompanies = "SELECT * FROM company"

// Company represents the company table structure
type Company struct {
	ID 		int
	Name 	string
	Type	string
}

var tmpl = template.Must(template.ParseGlob("forms/*"))

// ShowAllCompanies displays all of the companies
func ShowAllCompanies(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
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


