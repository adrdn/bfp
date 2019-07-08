package main

import (
	"log"
	"net/http"

	"drdn/bfp/company"
)

func main() {
	log.Println("Server is started on: http://localhost:8100")
	
	http.HandleFunc("/", company.ShowAllCompanies)

	http.ListenAndServe(":8100", nil)
}