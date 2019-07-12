package main

import (
	"log"
	"net/http"

	"adrdn/dit/user"
	"adrdn/dit/credential"
)

func main() {
	log.Println("Server is started on: http://localhost:8100")
	
	http.HandleFunc("/admin/users", user.DisplayAllUsers)
	// http.HandleFunc("/company/new", company.New)
	// http.HandleFunc("/company/insert", company.Insert)
	// //http.HandleFunc("/company/edit", company.Edit)
	// http.HandleFunc("/company/update", company.Update)
	http.HandleFunc("/admin/users/delete", user.DeleteUser)

	http.HandleFunc("/register", credential.SignUp)
	http.HandleFunc("/signup", credential.RegisterNewUser)
	http.HandleFunc("/login", credential.Login)
	http.HandleFunc("/auth", credential.Authentication)
	http.HandleFunc("/home", credential.Home)
	http.HandleFunc("/logout", credential.Logout)

	http.ListenAndServe(":8100", nil)
}