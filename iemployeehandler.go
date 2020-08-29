package handlers

import "net/http"

//EmployeeHandlerIface abstracts all the functionalities of the employees handler
type EmployeeHandlerIface interface {
	SearchEmployee(w http.ResponseWriter, r *http.Request)
	PostEmployee(w http.ResponseWriter, r *http.Request)
	ListEmployees(w http.ResponseWriter, r *http.Request)
	UpdateEmployeeDetails(w http.ResponseWriter, r *http.Request)
	RestoreEmployeeStatus(w http.ResponseWriter, r *http.Request)
	GetEmployee(w http.ResponseWriter, r *http.Request)
	DeleteEmployee(w http.ResponseWriter, r *http.Request)
}
