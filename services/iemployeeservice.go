package services

import "aicumen/models"

//EmployeeServiceIface abstracts the methods for EmployeeService Implementation
type EmployeeServiceIface interface {
	PostEmployee(employee models.Employee) (string, error)
	ListAllEmployees() ([]models.Employee, error)
	SearchEmployee(query string) ([]models.Employee, error)
	PatchEmployeeDetails(employeeID string, employeeDetails models.Employee) error
	GetEmployee(employeeID string) (models.Employee, error)
	DeleteEmployee(employeeID string) error
}
