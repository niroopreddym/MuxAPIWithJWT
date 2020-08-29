package enums

//EmployeeStatus is the enum for Employee Status
type EmployeeStatus int

const (
	//ACTIVE is one of the ServiceTypes
	ACTIVE EmployeeStatus = iota
	//INACTIVE is one of the ServiceTypes
	INACTIVE
)

//ArrEmployeeStatus represent all the registered EmployeeStatus enum descriptions
var ArrEmployeeStatus = []string{"ACTIVE", "INACTIVE"}

func (s EmployeeStatus) String() string {
	return ArrEmployeeStatus[s]
}
