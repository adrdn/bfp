package flow

import (
	"net/http"
	"text/template"

	"adrdn/dit/config"
)

const echoAllFlows	= "SELECT name FROM flow"

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

