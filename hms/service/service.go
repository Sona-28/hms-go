package service

import (
	"context"
	"errors"
	"fmt"
	"hms/config"
	"hms/models"
	"log"
	"math/rand"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllCustomers() ([]models.Customer, error) {

	filter := bson.D{}
	cursor, err := config.Customer_ProfileCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var customers []models.Customer
	for cursor.Next(context.Background()) {
		var customer models.Customer
		err := cursor.Decode(&customer)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		customers = append(customers, customer)
	}
	fmt.Println(customers)
	return customers, nil
}

func Insert(profile models.Customer) error {
	ctx := context.Background()

	// Check if a customer with the same email already exists.
	query := bson.M{"email": profile.Email}
	var existingCustomer models.Customer
	err := config.Customer_ProfileCollection.FindOne(ctx, query).Decode(&existingCustomer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No customer with the same email exists, so generate a new ID and insert the profile.
			id := GenerateID()
			fmt.Println(id)
			// check := checkID(id)
			// if check == false {}

			// Set the patient ID in the profile.
			profile.PatientID = id
			inserted, err := config.Customer_ProfileCollection.InsertOne(ctx, profile)
			if err != nil {
				log.Fatal(err)

				fmt.Println(err)
				return err

			}
			fmt.Println("Inserted", inserted.InsertedID)
			return nil
		}
		// Some other error occurred during the query.
		log.Fatal(err)
		return err
	}

	// A customer with the same email already exists, return an error.
	return err
}

// func checkID(id string) bool {
// 	var ctx context.Context

// 	var existingCustomer models.Customer
// 	cursor := bson.M{"id": id}
// 	err := config.Customer_ProfileCollection.FindOne(ctx, cursor).Decode(&existingCustomer)

// 	if err != nil {
// 		return true
// 	}
// 	return false
// }

func Appoitment(profile models.Appoitment) {
	inserted, err := config.Customer_Collection.InsertOne(context.Background(), profile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted", inserted.InsertedID)
}
func Login(profile models.Login) error {
	fmt.Println("service")
	var ctx context.Context
	//  var login models.Login
	query := bson.M{"email": profile.Email, "password": profile.Password}
	var customer models.Customer
	err := config.Customer_ProfileCollection.FindOne(ctx, query).Decode(&customer)
	if err != nil {
		fmt.Println(err)
	}
	return err

}

func Feedback(profile models.Feedback) {
	fmt.Println("service")
	inserted, err := config.Customer_feedback.InsertOne(context.Background(), profile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted", inserted.InsertedID)
}

func Adminlogin(profile models.Adminlogin) error {
	var ctx context.Context
	query := bson.M{"adminid": profile.AdminID, "password": profile.Password}
	var Admin models.Adminlogin
	err := config.Customer_Admin.FindOne(ctx, query).Decode(&Admin)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func Createadmin(profile models.Createadmin) {
	inserted, err := config.Customer_Admin.InsertOne(context.Background(), profile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted", inserted.InsertedID)
}

func GenerateID() string {
	length := 4
	// Define a character set from which the ID will be generated.
	characters := "0123456789"

	// Initialize a random seed based on the current time.
	rand.Seed(time.Now().UnixNano())

	// Create a byte slice to store the generated ID.
	id := make([]byte, length)

	// Generate a random ID by selecting random characters from the character set.
	for i := 0; i < length; i++ {
		id[i] = characters[rand.Intn(len(characters))]
	}

	// Convert the byte slice to a string and return the random ID.
	return "PA" + string(id)
}

func GetById(id string) (models.Customer, error) {
	var patient models.Customer
	filter := bson.M{"patientid": id}
	err := config.Customer_ProfileCollection.FindOne(context.Background(), filter).Decode(&patient)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle the case where no documents are found
			return patient, fmt.Errorf("No data found for the given ID: %s", id)
		}

		// Handle other errors
		return patient, err
	}
	return patient, nil
}

func DeleteById(adminId string) error {
	// Define a filter to find the document to be deleted
	filter := bson.M{"patientid": adminId}
	fmt.Println("service")
	var ctx context.Context

	var customer models.Customer // Replace with your actual data model
	if err := config.Customer_ProfileCollection.FindOneAndDelete(ctx, filter).Decode(&customer); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("No data found for the given Admin ID")
		}
		return err
	}

	return nil
}

func ViewAppointment(patientID string) ([]models.Appoitment, error) {
	//fmt.Println("service")
	// Define a filter to query appointments by patient ID.
	filter := bson.M{"patientid": patientID}

	var appointments []models.Appoitment
	cursor, err := config.Customer_Collection.Find(context.Background(), filter)
	for cursor.Next(context.Background()) {
		var appointment models.Appoitment
		if err := cursor.Decode(&appointment); err != nil {
			log.Fatal(err)
			return appointments, err
		}
		appointments = append(appointments, appointment)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle the case where no documents are found
			return appointments, fmt.Errorf("no data found for the given ID: %s", patientID)
		}

		// Handle other errors
		return appointments, err
	}
	// fmt.Println(appointments)
	return appointments, nil

}
func ViewAllAppointments() ([]models.Appoitment, error) {
	//fmt.Println("service")
	// Define a filter to query appointments by patient ID.
	filter := bson.M{}

	var appointments []models.Appoitment
	cursor, err := config.Customer_Collection.Find(context.Background(), filter)
	for cursor.Next(context.Background()) {
		var appointment models.Appoitment
		if err := cursor.Decode(&appointment); err != nil {
			log.Fatal(err)
			return appointments, err
		}
		appointments = append(appointments, appointment)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle the case where no documents are found
			return appointments, fmt.Errorf("no data found")
		}

		// Handle other errors
		return appointments, err
	}
	// fmt.Println(appointments)
	return appointments, nil

}
func ViewFeedback() ([]models.Feedback, error) {
	filter := bson.D{}
	cursor, err := config.Customer_feedback.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var customers []models.Feedback
	for cursor.Next(context.Background()) {
		var customer models.Feedback
		err := cursor.Decode(&customer)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		customers = append(customers, customer)
	}
	fmt.Println(customers)
	return customers, nil
}

func GetLoginDetails(id string) (models.Customer, error){
	var patient models.Customer
	filter := bson.M{"email": id}
	err := config.Customer_ProfileCollection.FindOne(context.Background(), filter).Decode(&patient)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle the case where no documents are found
			return patient, fmt.Errorf("no data found for the given ID: %s", id)
		}

		// Handle other errors
		return patient, err
	}
	return patient, nil
}

func PredictDisease() ([]models.Prediction, error) {
	appointment, err := ViewAllAppointments()
	if err != nil {
		return nil, err
	}
	fmt.Println("Predicting severity for the next 24 hours...")

	// Get current date
	currentDate := time.Now().Format("2006-01-02")

	// Count occurrences of diseases for the current date in historical data
	diseaseCounts := make(map[string]int)
	for _, data := range appointment {
		if data.Date == currentDate {
			diseaseCounts[data.Purpose]++
		}
	}

	// If no data is available for the current date, predict based on all historical data
	if len(diseaseCounts) == 0 {
		for _, data := range appointment {
			diseaseCounts[data.Purpose]++
		}
	}

	// Predicted severity is the total count of diseases for the current date
	totalOccurrences := len(appointment)
	fmt.Printf("Predicted severity for %s:\n", currentDate)
	var predict []models.Prediction
	for disease, count := range diseaseCounts {
		severity := (count * 10) / totalOccurrences // Scale severity from 0 to 10 based on frequency
		fmt.Printf("%s: %d\n", disease, severity)
		r := models.Prediction{
			DiseaseName: disease,
			Severity:    severity,
		}
		predict = append(predict, r)
	}
	return predict,nil
}
