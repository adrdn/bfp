package home

import (
	"fmt"

	"adrdn/dit/config"
)

const pending = "SELECT request_ID FROM pending WHERE role = ?"
const pendingCount = "SELECT COUNT(*) FROM pending WHERE role = ?"

func showPendingID(role string)[]string {
	db := config.DbConn()
	var ID 		string
	var IDList 	[]string

	res, err := db.Query(pending, role)
	if err != nil {
		fmt.Println(err)
	}
	for res.Next() {
		err = res.Scan(&ID)
		if err != nil {
			fmt.Println(err)
		}
		IDList = append(IDList, ID)
	}
	return IDList
}

func countPending(role string) int {
	db := config.DbConn()
	var count int

	res, err := db.Query(pendingCount, role)
	if err != nil {
		fmt.Println(err)
	}
	for res.Next() {
		err = res.Scan(&count)
		if err != nil {
			fmt.Println(err)
		}
	}
	return count
}
