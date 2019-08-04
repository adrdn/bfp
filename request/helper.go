package request

import(
	"fmt"
	"strings"

	"adrdn/dit/config"
)

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