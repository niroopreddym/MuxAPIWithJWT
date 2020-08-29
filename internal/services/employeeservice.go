package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	client "muxapiwithjwt/internal/clients"
	"muxapiwithjwt/internal/models"

	"github.com/aws/aws-sdk-go/aws"
	uuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

//DbConnectionString should be read from environment
const DbConnectionString = "mongodb://192.168.99.100:27017"

//DbName is the name of my DB in Mongo
const DbName = "aicumendb"

//CollectionName is the table in mongoDB
const CollectionName = "employee"

//EmployeeService provides the implementation of EmployeeServiceIface
type EmployeeService struct {
	mongoClient *mongo.Client
}

//NewEmployeeService is the constructor for EmployeeService Struct
func NewEmployeeService() *EmployeeService {
	newClient, err := client.NewMongoClient(DbConnectionString)
	if err != nil {
		fmt.Println("some error occured while creating mongo client ", err.Error())
	}

	dbInstance := newClient.Database(DbName)
	collection := dbInstance.Collection(CollectionName)

	indexName, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{
					Key:   "employeeid",
					Value: bsonx.Int32(1),
				},
			},
			Options: options.Index().SetUnique(true),
		},
	)

	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Options: options.Index().SetName("TextIndex"),
			Keys: bsonx.Doc{
				{
					Key:   "$**",
					Value: bsonx.String("text"),
				},
			},
		},
	)

	fmt.Println(indexName)
	return &EmployeeService{
		mongoClient: newClient,
	}
}

//PostEmployee adds an employee to the DB
func (service *EmployeeService) PostEmployee(employee models.Employee) (string, error) {
	UUID := uuid.New().String()
	employee.EmployeeID = &UUID
	employee.Status = aws.String("ACTIVE")
	collection := service.mongoClient.Database(DbName).Collection(CollectionName)
	result, err := collection.InsertOne(context.Background(), employee)
	fmt.Println(result)
	if err != nil {
		return "", err
	}

	return UUID, nil
}

//ListAllEmployees list down all the employees from the DB
func (service *EmployeeService) ListAllEmployees() ([]models.Employee, error) {
	collection := service.mongoClient.Database(DbName).Collection(CollectionName)

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var employees []models.Employee
	if err = cursor.All(context.Background(), &employees); err != nil {
		log.Fatal(err)
	}

	return employees, nil
}

//SearchEmployee lists down all the employees who match a search query
func (service *EmployeeService) SearchEmployee(query string) ([]models.Employee, error) {
	collection := service.mongoClient.Database(DbName).Collection(CollectionName)

	cursor, err := collection.Find(context.Background(), bson.M{
		"$text": bson.M{
			"$search": query,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	var employees []models.Employee
	if err = cursor.All(context.Background(), &employees); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return employees, nil
}

//PatchEmployeeDetails updates the employee details
func (service *EmployeeService) PatchEmployeeDetails(employeeID string, employeeDetails models.Employee) error {
	collection := service.mongoClient.Database(DbName).Collection(CollectionName)
	updatesToBePerformed := bson.M{}
	updatesToBePerformed["employeeid"] = employeeID
	if employeeDetails.Department != nil {
		updatesToBePerformed["department"] = employeeDetails.Department
	}

	if employeeDetails.Name != nil {
		updatesToBePerformed["name"] = employeeDetails.Name
	}

	if employeeDetails.Skills != nil {
		updatesToBePerformed["skills"] = employeeDetails.Skills
	}

	if employeeDetails.Address != nil {
		address := models.Address{}
		if employeeDetails.Address.City != nil {
			address.City = employeeDetails.Address.City
		}

		if employeeDetails.Address.Country != nil {
			address.Country = employeeDetails.Address.Country
		}

		if employeeDetails.Address.DoorNo != nil {
			address.DoorNo = employeeDetails.Address.DoorNo
		}

		if employeeDetails.Address.State != nil {
			address.State = employeeDetails.Address.State
		}

		updatesToBePerformed["address"] = address
	}

	if employeeDetails.Status != nil {
		updatesToBePerformed["status"] = employeeDetails.Status
	}

	// consolidatedMap(&updatesToBePerformed, employeeDetails)

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"employeeid": employeeID},
		bson.M{
			"$set": updatesToBePerformed,
		})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)

	return nil
}

//GetEmployee gets the employee details with the given employeeID
func (service *EmployeeService) GetEmployee(employeeID string) (models.Employee, error) {
	collection := service.mongoClient.Database(DbName).Collection(CollectionName)

	cursor, err := collection.Find(context.Background(), bson.M{"employeeid": string(employeeID)})
	if err != nil {
		log.Fatal(err)
	}

	var employees []models.Employee
	if err = cursor.All(context.Background(), &employees); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	if len(employees) > 0 {
		return employees[0], nil
	}

	return models.Employee{}, errors.New("No User Found")
}

//DeleteEmployee deletes an employee permenantly from collection
func (service *EmployeeService) DeleteEmployee(employeeID string) error {
	collection := service.mongoClient.Database(DbName).Collection(CollectionName)
	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"employeeid": employeeID})
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(deleteResult.DeletedCount)

	return nil
}
