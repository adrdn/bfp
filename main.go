package main

import (
	"log"
	"net/http"

	"adrdn/bfp/company"
	"adrdn/bfp/credential"
)

func main() {
	log.Println("Server is started on: http://localhost:8100")
	
	http.HandleFunc("/company", company.ShowAllCompanies)
	http.HandleFunc("/company/new", company.New)
	http.HandleFunc("/company/insert", company.Insert)
	http.HandleFunc("/company/edit", company.Edit)
	http.HandleFunc("/company/update", company.Update)
	http.HandleFunc("/company/delete", company.Delete)

	http.HandleFunc("/register", credential.SignUp)
	http.HandleFunc("/login", credential.RegisterNewUser)

	http.ListenAndServe(":8100", nil)
}