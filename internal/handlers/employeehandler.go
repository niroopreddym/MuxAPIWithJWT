package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"github.com/niroopreddym/muxapiwithjwt/internal/enums"
	"github.com/niroopreddym/muxapiwithjwt/internal/models"
	"github.com/niroopreddym/muxapiwithjwt/internal/services"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gorilla/mux"
)

//EmployeeHandler handles all the operations related to the Employees
type EmployeeHandler struct {
	EmployeeService services.EmployeeServiceIface
}

//NewEmployeeHandler returns a new instance of the EmployeeHandler
func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{
		EmployeeService: services.NewEmployeeService(),
	}
}

//PostEmployee handler handles the bussiness logic for posting a new Employee
func (handler *EmployeeHandler) PostEmployee(w http.ResponseWriter, r *http.Request) {
	employeeDetails := models.Employee{}
	bodyBytes, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		responseController(w, http.StatusInternalServerError, readErr)
		return
	}

	strBufferValue := string(bodyBytes)
	err := json.Unmarshal([]byte(strBufferValue), &employeeDetails)
	if err != nil {
		responseController(w, http.StatusInternalServerError, err)
		return
	}

	errorMessages := []string{}
	postRequestBodyInitialValidation(employeeDetails, &errorMessages)
	if len(errorMessages) > 0 {
		responseController(w, http.StatusBadRequest, errorMessages)
		return
	}

	uniqueID, err := handler.EmployeeService.PostEmployee(employeeDetails)
	if err != nil {
		fmt.Println(err)
		responseController(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseController(w, http.StatusOK, map[string]string{
		"employeeId": uniqueID,
	})
}

//GetEmployee gets specific employee details
func (handler *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID := params["employeeId"]

	employeeDetails, err := handler.EmployeeService.GetEmployee(employeeID)
	if err != nil {
		fmt.Println(err)
		responseController(w, http.StatusInternalServerError, "Error occured while fetching the employee details")
		return
	}

	responseController(w, http.StatusOK, employeeDetails)
}

//ListEmployees handler handles the bussiness logic for listing all the Employees
func (handler *EmployeeHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	lstEmployees, err := handler.EmployeeService.ListAllEmployees()
	if err != nil {
		fmt.Println(err)
		responseController(w, http.StatusInternalServerError, "Error occured while fetch the userdetails")
		return
	}

	responseController(w, http.StatusOK, lstEmployees)
}

//UpdateEmployeeDetails handles the bussiness logic for Updating/deleting the EmployeeDetails
func (handler *EmployeeHandler) UpdateEmployeeDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID := params["employeeId"]
	employeeDetails, err := handler.EmployeeService.GetEmployee(employeeID)
	if err != nil {
		responseController(w, http.StatusInternalServerError, err.Error())
		return
	}

	if *employeeDetails.EmployeeID == "" {
		responseController(w, http.StatusBadRequest, "Employee Not Found")
		return
	}

	requestedEmployeeDetails := models.Employee{}
	bodyBytes, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		fmt.Println("reading from buffer")
		err := errors.New("error reading data from the response " + readErr.Error())
		fmt.Println(err)
		responseController(w, http.StatusInternalServerError, err.Error())
		return
	}

	strBufferValue := string(bodyBytes)
	err = json.Unmarshal([]byte(strBufferValue), &requestedEmployeeDetails)
	if err != nil {
		fmt.Println(err)
		responseController(w, http.StatusInternalServerError, err.Error())
		return
	}

	errorMessages := []string{}
	patchCallInitialValidation(requestedEmployeeDetails, &errorMessages)
	if len(errorMessages) > 0 {
		responseController(w, http.StatusBadRequest, errorMessages)
		return
	}

	err = handler.EmployeeService.PatchEmployeeDetails(employeeID, requestedEmployeeDetails)
	if err != nil {
		fmt.Println(err)
		responseController(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseController(w, http.StatusNoContent, "Updated Sucessfully")
}

//RestoreEmployeeStatus handles the bussiness logic for Restoring Soft deleted employee
func (handler *EmployeeHandler) RestoreEmployeeStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID := params["employeeId"]

	employeeDetails, err := handler.EmployeeService.GetEmployee(employeeID)
	if err != nil {
		fmt.Println(err)
		responseController(w, http.StatusInternalServerError, "Error occured while fetching the employee details")
		return
	}

	if *employeeDetails.Status == enums.EmployeeStatus(enums.ACTIVE).String() {
		responseController(w, http.StatusBadRequest, "Cannot Reactivate an Active Employee")
		return
	}

	requestedEmployeeDetails := models.Employee{
		Status: aws.String(enums.EmployeeStatus(enums.ACTIVE).String()),
	}

	err = handler.EmployeeService.PatchEmployeeDetails(employeeID, requestedEmployeeDetails)
	if err != nil {
		responseController(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseController(w, http.StatusNoContent, "Re-Activated Sucessfully")
}

//SearchEmployee handles the bussiness logic for retrieving the list of employees based on a search query
func (handler *EmployeeHandler) SearchEmployee(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var qpValue string

	if len(queryParams) == 0 {
		responseController(w, http.StatusBadRequest, "Invalid Query Param: Please pass any qualified string")
		return
	}

	for qpKey, values := range queryParams {
		if qpKey != "query" {
			responseController(w, http.StatusBadRequest, "Invalid Query Param")
			return
		}

		if len(values) > 1 {
			responseController(w, http.StatusBadRequest, "only one Query Param value is allowed")
			return
		}

		if values[0] == "" {
			responseController(w, http.StatusBadRequest, "Invalid Query Param: Please pass any qualified string")
			return
		}

		qpValue = values[0]
	}

	lstEmployees, err := handler.EmployeeService.SearchEmployee(qpValue)
	if err != nil {
		responseController(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseController(w, http.StatusOK, lstEmployees)
}

//DeleteEmployee deletes/inactivates an employee
func (handler *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	employeeID := params["employeeId"]

	queryParams := r.URL.Query()
	permanentDeleteFlag := false
	for qpKey := range queryParams {
		if qpKey == "permanentlyDelete" {
			permanentDeleteFlag = true
		} else {
			responseController(w, http.StatusBadRequest, "Invalid Query Param value")
		}
	}

	if permanentDeleteFlag {
		err := handler.EmployeeService.DeleteEmployee(employeeID)
		if err != nil {
			responseController(w, http.StatusInternalServerError, err.Error())
		}

		responseController(w, http.StatusNoContent, "Successfully Deleted")
		return
	}

	requestedEmployeeDetails := models.Employee{
		Status: aws.String(enums.EmployeeStatus(enums.INACTIVE).String()),
	}

	err := handler.EmployeeService.PatchEmployeeDetails(employeeID, requestedEmployeeDetails)
	if err != nil {
		responseController(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseController(w, http.StatusNoContent, "Deleted Sucessfully")
}

func responseController(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func postRequestBodyInitialValidation(employeeDetails models.Employee, errorMessages *[]string) {
	if employeeDetails.Name == nil {
		errorMessage := "Attribute Missing: Name in the request body"
		*errorMessages = append(*errorMessages, errorMessage)
	}
}

func patchCallInitialValidation(employeeDetails models.Employee, errorMessages *[]string) {
	if employeeDetails.Status != nil {
		errorMessage := "Invalid Attribute: Status in the request body"
		*errorMessages = append(*errorMessages, errorMessage)
	}
}
