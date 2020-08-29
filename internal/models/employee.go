package models

//Employee model exposes the required attributes for the Employee table
type Employee struct {
	EmployeeID *string  `json:"employeeId,omitempty"`
	Name       *string  `json:"name"`
	Department *string  `json:"department,omitempty"`
	Address    *Address `json:"address,omitempty"`
	Skills     []string `json:"skills,omitempty"`
	Status     *string  `json:"status,omitempty"`
}
