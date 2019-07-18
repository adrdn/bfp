package workflow

import (
	"net/http"
	"text/template"

	"adrdn/dit/flow"
	"adrdn/dit/config"
)

const echoAllFlow = "SELECT * FROM flow"
const addNewRequest = "INSERT INTO request(type, current_step) VALUES(?, ?)"

var tmpl = template.Must(template.ParseGlob("forms/workflow/*"))

// Flow defines the selected flow by the user
var Flow string
// Description defines the description entered by the user
var Description string

// New starts a new request of the select flow
func New(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	flowEntity := flow.Flow{}
	flowList := []flow.Flow{}
	
	flows, err := db.Query(echoAllFlow)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	for flows.Next() {
		err = flows.Scan(&flowEntity.ID, &flowEntity.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		flowList = append(flowList, flowEntity)
	}
	db.Close()
	tmpl.ExecuteTemplate(w, "New", flowList)
}

// Insert adds the new entity
func Insert(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	if r.Method == "POST" {
		selectedFlow := r.FormValue("flow")
		desc := r.FormValue("description")

		newRequest, err :=db.Prepare(addNewRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// current_step is assinged to 2 because at this stage the request is already created
		_, err = newRequest.Exec(selectedFlow, 2)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		Description = desc
	}
	db.Close()
	http.Redirect(w, r, "/request/view", 301)
}

// Echo displays all of the requests
func Echo(w http.ResponseWriter, r *http.Request) {

}