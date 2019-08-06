package request

import(
	"fmt"
	"strings"

	"adrdn/dit/config"
)

const echoAllFlow = "SELECT * FROM flow"
const echoALLRequest = "SELECT * FROM request"
const echoOneRequest = "SELECT ID, type, current_step, termination, completion, description FROM request WHERE ID = ?"
const addNewRequest = "INSERT INTO request(type, current_step, termination, completion, deletion, description, created_at, created_by, updated_at, updated_by) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const updateRequest = "UPDATE request SET current_step = ?, description = ?, updated_at = ?, updated_by = ? WHERE ID = ?"
const addNewPending = "INSERT INTO pending(request_ID, role) VALUES (?, ?)"
const updatePending = "UPDATE pending SET role = ? WHERE request_ID = ?"
const terminateRequest = "UPDATE request SET current_step = 0, termination = ?, description = ?, updated_at = ?, updated_by = ? WHERE ID = ?"
const finishRequest = "UPDATE request SET current_step = 0, completion = ?, description = ?, updated_at = ?, updated_by = ? WHERE ID = ?"
const fetchTotalSteps = "SELECT total_steps from flow_"
const deleteRequest = "UPDATE request SET deletion = ? WHERE ID = ?"
// const populateUpdatedAt = "UPDATE request SET updated_at = ? WHERE ID = ?"
// const populateUpdatedBy = "UPDATE request SET updated_by = ? WHERE ID = ?"

const terminatedStatus 	=	"Terminated"
const completedStatus 	=	"Completed"
const runningStatus		=	"In Process"

func getStepValue(stepNumber, flowName string) string {
	db := config.DbConn()

	var step string
	stepValue, err := db.Query("SELECT step" + stepNumber + " FROM flow_" + strings.ToUpper(flowName))
	if err != nil {
		fmt.Println(err)
	}
	for stepValue.Next() {
		err = stepValue.Scan(&step)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer db.Close()
	return step
}

func addPending(requestID int, roleName string) {
	db := config.DbConn()
	
	entity, err := db.Prepare(addNewPending)
	if err != nil {
		fmt.Println(err)
	}
	_, err = entity.Exec(requestID, roleName)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}

func updatePendingTable(role string, requestID int) {
	db := config.DbConn()

	res, err := db.Prepare(updatePending)
	if err != nil {
		fmt.Println(err)
	}
	res.Exec(role, requestID)
	defer db.Close()
}

// func updatedAt(ID string) {
// 	db := config.DbConn()

// 	t := time.Now()
// 	dateTimeLayout := t.Format("2006-01-02 15:04:05")

// 	res, err := db.Prepare(populateUpdatedAt)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	res.Exec(dateTimeLayout, ID)
// 	defer db.Close()
// }

// func updatedBy(username, ID string) {
// 	db := config.DbConn()

// 	res, err := db.Prepare(populateUpdatedBy)
// 	if err != nil {
// 		println(err)
// 	}
// 	res.Exec(username, ID)
// 	defer db.Close()
// }