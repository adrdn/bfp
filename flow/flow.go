package flow

import (
	"net/http"
	"text/template"

	"adrdn/dit/config"
)

const echoAllFlows	= "SELECT * FROM flow"
const newFlow		= "INSERT INTO flow (name) VALUES (?)"
const deleteFlow	= "Delete FROM flow WHERE id = ?"

// Flow represents the flow structure
type Flow struct {
	ID int
	Name string
}

var tmpl = template.Must(template.ParseGlob("forms/admin/flow/*"))

// ShowAllFlows displays all of the roles
func ShowAllFlows (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	rList, err := db.Query(echoAllFlows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	flow := Flow{}
	flowList := []Flow{}

	for rList.Next() {
		err = rList.Scan(&flow.ID, &flow.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		flowList = append(flowList, flow)
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Echo", flowList)
}

// New revokes the new page
func New (w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", 301)
}

// Insert adds the new flow
func Insert (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		insForm, err := db.Prepare(newFlow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		insForm.Exec(name)
	}
	defer db.Close()
	http.Redirect(w, r, "/admin/flow", 301)
}

// Delete removes the flow
func Delete (w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	id := r.URL.Query().Get("id")
	delForm, err := db.Prepare(deleteFlow)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	delForm.Exec(id)
	defer db.Close()
	http.Redirect(w, r, "/admin/flow", 301)
}

