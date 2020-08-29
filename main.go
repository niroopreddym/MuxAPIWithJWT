package main

import (
	"muxapiwithjwt/handlers"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func main() {
	//key could be saved somewhere safe
	var mySigningKey = []byte("ultimateStarAjith")

	jwtMiddlewareInstance := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter()
	employeeHandler := handlers.NewEmployeeHandler()
	authHandler := handlers.NewAuthHandler()
	fmt.Println("Aicumen employee server started listening on port : ", 9293)

	router.Handle("/employees/search", jwtMiddlewareInstance.Handler(http.HandlerFunc(employeeHandler.SearchEmployee))).Methods("GET")
	router.Handle("/employees/add", jwtMiddlewareInstance.Handler(http.HandlerFunc(employeeHandler.PostEmployee))).Methods("POST")
	router.Handle("/employees", jwtMiddlewareInstance.Handler(http.HandlerFunc(employeeHandler.ListEmployees))).Methods("GET")
	router.Handle("/employees/{employeeId}", jwtMiddlewareInstance.Handler(http.HandlerFunc(employeeHandler.GetEmployee))).Methods("GET")
	router.Handle("/employees/{employeeId}", jwtMiddlewareInstance.Handler(http.HandlerFunc(employeeHandler.UpdateEmployeeDetails))).Methods("PATCH")
	router.Handle("/employees/{employeeId}", jwtMiddlewareInstance.Handler(http.HandlerFunc(employeeHandler.DeleteEmployee))).Methods("DELETE")
	router.Handle("/employees/{employeeId}/restore", jwtMiddlewareInstance.Handler(http.HandlerFunc(employeeHandler.RestoreEmployeeStatus))).Methods("PATCH")
	router.HandleFunc("/getjwttoken", authHandler.GenerateToken).Methods("GET")

	http.ListenAndServe(":9293", router)
}
