package models

//Address address of the employee
type Address struct {
	DoorNo  *string `json:"doorNo,omitempty"`
	City    *string `json:"city,omitempty"`
	State   *string `json:"state,omitempty"`
	Country *string `json:"country,omitempty"`
}
