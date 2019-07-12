package main

import (
	"log"
	"net/http"

	"adrdn/dit/user"
	"adrdn/dit/role"
	"adrdn/dit/credential"
)

func main() {
	log.Println("Server is started on: http://localhost:8000")
	
	http.HandleFunc("/admin/users", user.DisplayAllUsers)
	http.HandleFunc("/admin/users/delete", user.DeleteUser)

	http.HandleFunc("/admin/role", role.ShowAllRoles)
	http.HandleFunc("/admin/edit", role.Edit)
	http.HandleFunc("/admin/role/update", role.Update)

	http.HandleFunc("/register", credential.SignUp)
	http.HandleFunc("/signup", credential.RegisterNewUser)
	http.HandleFunc("/login", credential.Login)
	http.HandleFunc("/auth", credential.Authentication)
	http.HandleFunc("/home", credential.Home)
	http.HandleFunc("/logout", credential.Logout)

	http.ListenAndServe(":8000", nil)
}