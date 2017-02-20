package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/tables", tables)
	http.HandleFunc("/person", person)
	http.HandleFunc("/party", party)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!\n"))
}

func tables(w http.ResponseWriter, r *http.Request) {
	data, err := getTables()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
	//w.Write([]byte(data))
}

func person(w http.ResponseWriter, r *http.Request) {
	data, err := getPerson()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var person PersonData
	jsonErr := dynamodbattribute.UnmarshalMap(data, &person)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(person)
	//w.Write([]byte(person.Name))
}

func party(w http.ResponseWriter, r *http.Request) {
	data, err := getParty()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var party PartyData
	jsonErr := dynamodbattribute.UnmarshalMap(data, &party)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(party)
	//w.Write([]byte(person.Name))
}

func getTables() (string, error) {
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2")}))
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
	    return "", err
	}

	var tables string = ""
	for _, table := range result.TableNames {
	    tables = tables + *table
	}
	return tables, nil
}

func getPerson() (map[string]*dynamodb.AttributeValue, error) {
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config *aws.Config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("err")
	}
	var svc *dynamodb.DynamoDB = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	getItemInput.SetTableName("Person")
	var attributeValue = dynamodb.AttributeValue{}
	attributeValue.SetS("12345699033")
	getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"facebookID": &attributeValue})
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	return getItemOutput.Item, err2
}

func getParty() (map[string]*dynamodb.AttributeValue, error) {
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config *aws.Config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("err")
	}
	var svc *dynamodb.DynamoDB = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	getItemInput.SetTableName("Party")
	var attributeValue = dynamodb.AttributeValue{}
	attributeValue.SetN("1")
	getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"partyID": &attributeValue})
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	return getItemOutput.Item, err2
}

type PersonData struct {
		BlackoutList []string `json:"blackoutList"`
		FacebookID string 		`json:"facebookID"`
		InvitedTo []uint64 		`json:"invitedTo"`
		IsABarOwner bool 			`json:"isABarOwner"`
		IsMale bool 					`json:"isMale"`
		Name string 					`json:"name"`
		PartyHostFor []uint64 `json:"partyHostFor"`
}

type PartyData struct {
	AddressLine1        string  	`json:"addressLine1"`
	AddressLine2        string  	`json:"addressLine2"`
	City                string  	`json:"city"`
	Country             string  	`json:"country"`
	Details	         		string  	`json:"details"`
	DrinksProvided      bool    	`json:"drinksProvided"`
	EndTime             string  	`json:"endTime"`
	FeeForDrinks        bool    	`json:"feeForDrinks"`
	Hosts								[]Host 		`json:"hosts"`
	Invitees						[]Invitee `json:"invitees"`
	Latitude            float64 	`json:"latitude"`
	Longitude           float64 	`json:"longitude"`
	PartyID             uint64   	`json:"partyID"`
	StartTime           string  	`json:"startTime"`
	StateProvinceRegion string  	`json:"stateProvinceRegion"`
	Title           		string  	`json:"title"`
	ZipCode             uint32  	`json:"zipCode"`
}

type Host struct {
	FacebookID	string 	`json:"facebookID"`
	IsMainHost 	bool 		`json:"isMainHost"`
	Name				string 	`json:"name"`
}

type Invitee struct {
	AtParty 		bool		`json:"atParty"`
	FacebookID	string	`json:"facebookID"`
	IsMale			bool		`json:"isMale"`
	Name				string	`json:"name"`
	Rating			string	`json:"rating"`
	Status			string 	`json:"status"`
}
