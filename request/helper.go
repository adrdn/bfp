package request

import(
	"fmt"
	"strings"

	"adrdn/dit/config"
)

func getType(ID int) string {
	db := config.DbConn()
	var typeValue string
	readType, err := db.Query(fetchType, ID)
	if err != nil {
		fmt.Println(err)
	}
	for readType.Next() {
		err = readType.Scan(&typeValue)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer db.Close()
	return typeValue
}

func getStep(stepNumber string, flowType string) string {
	db := config.DbConn()
	var stepValue string
	readStep, err := db.Query("SELECT step" + stepNumber + " FROM flow_" + strings.ToUpper(flowType))
	if err != nil {
		fmt.Println(err)
	}
	for readStep.Next() {
		err = readStep.Scan(&stepValue)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer db.Close()
	return stepValue
}

func getPending(roleName string) string {
	db := config.DbConn() 
	var pending string
	readPending, err := db.Query(fetchPendingValue, roleName)
	if err != nil {
		fmt.Println(err)
	}
	for readPending.Next() {
		err = readPending.Scan(&pending)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer db.Close()
	return pending
}

func updatePending(pendingValue, reqID, roleName string) {
	db := config.DbConn()
	res, err := db.Prepare(updateRolePending)
	if err != nil {
		fmt.Println(err)
	}

	if pendingValue == "" {
		_, err = res.Exec(reqID, roleName)
		if err != nil {
			fmt.Println(err)
		} 
	} else {
		_, err = res.Exec(pendingValue + ", " + reqID, roleName)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer db.Close()
}