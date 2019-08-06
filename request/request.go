package request

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"net/http"
	"text/template"

	"adrdn/dit/flow"
	"adrdn/dit/config"
	"adrdn/dit/credential"
)

var tmpl = template.Must(template.ParseGlob("forms/request/*"))

// Request represents the request structure
type Request struct {
	ID	 			int
	Type 			string
	PriorStep		string
	CurrentStep 	int
	StrCurrentStep	string
	NextStep		string
	Termination 	string
	Completion		string
	Deletion		string
	CreatedAt		string
	CreatedBy		string
	UpdatedAt		string
	UpdatedBy		string
	Status			string
	Description		string
	IsFirstStep		bool
	IsLastStep		bool
	TotalSteps		int
	IsDeleted		bool
}

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
	request := Request{}
	ok, user := credential.CheckAuthentication(w, r)
	if !ok {
		return
	}
	db := config.DbConn()
	if r.Method == "POST" {
		request.Type = r.FormValue("flow")
		request.Description = r.FormValue("description")
		request.CreatedBy = user.Username
		t := time.Now()
		request.CreatedAt = t.Format("2006-01-02 15:04:05")

		newRequest, err := db.Prepare(addNewRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// current_step is assinged to 2 because at this stage the request is already created
		res, err := newRequest.Exec(request.Type, 2, "", "", "", request.Description, request.CreatedAt, request.CreatedBy, "", "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// insert the ID and current step's role into the pending table
		stepValue := getStepValue("2", request.Type)
		ID, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		addPending(int(ID), stepValue)
	}
	db.Close()
	http.Redirect(w, r, "/request/view", 301)
}

// Echo displays all of the requests
func Echo(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	req := Request{}
	reqList := []Request{}

	requests, err := db.Query(echoALLRequest)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	for requests.Next() {
		err = requests.Scan(&req.ID, &req.Type, &req.CurrentStep, &req.Termination, &req.Completion, &req.Deletion, &req.Description, &req.CreatedAt, &req.CreatedBy, &req.UpdatedAt, &req.UpdatedBy)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if req.Deletion != "" {
			req.IsDeleted = true
		} else {
			req.IsDeleted = false
		}

		if req.Termination == "" && req.Completion == "" {
			req.Status = runningStatus
		} else if req.Termination != "" {
			req.Status = terminatedStatus	
		} else {
			req.Status = completedStatus
		}
		reqList = append(reqList, req)
	}

	db.Close()
	tmpl.ExecuteTemplate(w, "Echo", reqList)
}

// ShowDetails revoke the detail page
func ShowDetails(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	req := Request{}

	ID := r.URL.Query().Get("id")

	request, err := db.Query(echoOneRequest, ID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	for request.Next() {
		err = request.Scan(&req.ID, &req.Type, &req.CurrentStep, &req.Termination, &req.Completion, &req.Description)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if req.Termination != "" || req.Completion != "" {
		if req.Termination != "" {
			req.Status = terminatedStatus
		} else {
			req.Status = completedStatus
		}
		defer db.Close()
		tmpl.ExecuteTemplate(w, "Finished", req)
	} else {
		totalSteps, err := db.Query(fetchTotalSteps + strings.ToUpper(req.Type))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		for totalSteps.Next() {
			err = totalSteps.Scan(&req.TotalSteps)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		currentStep, err := db.Query("SELECT step" + strconv.Itoa(req.CurrentStep) + " FROM flow_" + strings.ToUpper(req.Type))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		for currentStep.Next() {
			err = currentStep.Scan(&req.StrCurrentStep)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		// If IsFirstStep is true, then show the terminate option on the Detail page
		if req.CurrentStep == 1 {
			req.IsFirstStep = true
			req.PriorStep = "You are the creator of this request"
		} else {
			// Fetch string value of the previous Step
			intPreStep := req.CurrentStep - 1
			stringPreStep := strconv.Itoa(intPreStep)	
			preStep, err := db.Query("SELECT step" + stringPreStep + " FROM flow_" + strings.ToUpper(req.Type))
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			for preStep.Next() {
				err = preStep.Scan(&req.PriorStep)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}

		// If IsLastStep is true, then show the Finalize option on the Detail page
		if req.CurrentStep == req.TotalSteps {
			req.IsLastStep = true
			req.NextStep = "You have to make the final decision"
		} else {
			// Fetch string value of the next Step
			intNextStep := req.CurrentStep + 1
			stringNextStep := strconv.Itoa(intNextStep)
			nextStep, err := db.Query("SELECT step" + stringNextStep + " FROM flow_" + strings.ToUpper(req.Type))
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			for nextStep.Next() {
				err = nextStep.Scan(&req.NextStep)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
		db.Close()
		tmpl.ExecuteTemplate(w, "Detail", req)
	}
}

// Update changes the request based on user decision
func Update(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	var stepValue string

	if r.Method == "POST" {
		ID := r.URL.Query().Get("id")
		intID, _ := strconv.Atoi(ID)
		currentStep := r.URL.Query().Get("cs")
		flowName := r.URL.Query().Get("type")
		intCurrentStep, _ := strconv.Atoi(currentStep)
		description := r.FormValue("description")
		decision := r.FormValue("decision")
		if decision == "approve" || decision == "reject" {
			if decision == "approve" {
				intCurrentStep++
				currentStep = strconv.Itoa(intCurrentStep)
			} else if decision == "reject" {
				intCurrentStep--
				currentStep = strconv.Itoa(intCurrentStep)
			}
			// update pending table
			stepValue = getStepValue(currentStep, flowName)
			updatePendingTable(stepValue, intID)
			// update request table
			request, err := db.Prepare(updateRequest)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			_, err = request.Exec(currentStep, description, ID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else if decision == "terminate" {
			request, err := db.Prepare(terminateRequest)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			t := time.Now()
			dateTimeLayout := t.Format("2006-01-02 15:04:05")
			_, err = request.Exec(dateTimeLayout, description, ID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			// update pending table
			updatePendingTable("", intID)
		} else {
			request, err := db.Prepare(finishRequest)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			t := time.Now()
			dateTimeLatout := t.Format("2006-01-02 15:04:05")
			_, err = request.Exec(dateTimeLatout, description, ID)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			// update pending table
			updatePendingTable("", intID)
		}
	}
	defer db.Close()
	http.Redirect(w, r, "/request/view", 301)
}

// Delete soft-deletes a request
func Delete(w http.ResponseWriter, r *http.Request) {
	db := config.DbConn()
	ID := r.URL.Query().Get("id")
	t := time.Now().Format("2006-01-02 15:04:05")
	request, err := db.Prepare(deleteRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = request.Exec(t, ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/request/view", 301)
}