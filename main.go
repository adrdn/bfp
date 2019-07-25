package main

import (
	"log"
	"net/http"

	"adrdn/dit/user"
	"adrdn/dit/role"
	"adrdn/dit/flow"
	"adrdn/dit/request"
	"adrdn/dit/credential"
)

func main() {
	log.Println("Server is started on: http://localhost:8000")

	http.HandleFunc("/register", credential.SignUp)
	http.HandleFunc("/signup", credential.RegisterNewUser)
	http.HandleFunc("/login", credential.Login)
	http.HandleFunc("/auth", credential.Authentication)
	http.HandleFunc("/home", credential.Home)
	http.HandleFunc("/refresh", credential.Refresh)
	
	http.HandleFunc("/admin/users", user.DisplayAllUsers)
	http.HandleFunc("/admin/users/delete", user.DeleteUser)

	http.HandleFunc("/admin/role", role.ShowAllRoles)
	http.HandleFunc("/admin/role/edit", role.Edit)
	http.HandleFunc("/admin/role/update", role.Update)
	http.HandleFunc("/admin/role/new", role.New)
	http.HandleFunc("/admin/role/insert", role.Insert)
	http.HandleFunc("/admin/role/delete", role.Delete)

	http.HandleFunc("/admin/flow", flow.ShowAllFlows)
	http.HandleFunc("/admin/flow/new", flow.New)
	http.HandleFunc("/admin/flow/insert", flow.Insert)
	http.HandleFunc("/admin/flow/delete", flow.Delete)

	http.HandleFunc("/request/new", request.New)
	http.HandleFunc("/request/process", request.Insert)
	http.HandleFunc("/request/view", request.Echo)
	http.HandleFunc("/request/detail", request.ShowDetails)
	http.HandleFunc("/request/update", request.Update)
	http.HandleFunc("/request/delete", request.Delete)

	http.ListenAndServe(":8000", nil)
}