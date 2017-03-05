package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

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
}

func bar(w http.ResponseWriter, r *http.Request) {
	data, err := getBar()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var bar BarData
	jsonErr := dynamodbattribute.UnmarshalMap(data, &bar)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(bar)
}

func getTables() (string, error) {
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2")}))
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return "", err
	}

	var tables string
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
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("err")
	}
	var svc = dynamodb.New(sess)
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
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("err")
	}
	var svc = dynamodb.New(sess)
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

func getBar() (map[string]*dynamodb.AttributeValue, error) {
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("err")
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	getItemInput.SetTableName("Bar")
	var attributeValue = dynamodb.AttributeValue{}
	attributeValue.SetN("1")
	getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"barID": &attributeValue})
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	return getItemOutput.Item, err2
}

// PersonData : A person from the database in json format
type PersonData struct {
	PeopleBlockingTheirActivityFromMe map[string]bool `json:"peopleBlockingTheirActivityFromMe"`
	FacebookID                        string          `json:"facebookID"`
	InvitedTo                         map[string]bool `json:"invitedTo"`
	IsMale                            bool            `json:"isMale"`
	Name                              string          `json:"name"`
	PartyHostFor                      map[string]bool `json:"partyHostFor"`
	BarHostFor                        map[string]bool `json:"barHostFor"`
}

// PartyData : A party from the database in json format
type PartyData struct {
	AddressLine1        string              `json:"addressLine1"`
	AddressLine2        string              `json:"addressLine2"`
	City                string              `json:"city"`
	Country             string              `json:"country"`
	Details             string              `json:"details"`
	DrinksProvided      bool                `json:"drinksProvided"`
	EndTime             string              `json:"endTime"`
	FeeForDrinks        bool                `json:"feeForDrinks"`
	Hosts               map[string]*Host    `json:"hosts"`
	Invitees            map[string]*Invitee `json:"invitees"`
	Latitude            float64             `json:"latitude"`
	Longitude           float64             `json:"longitude"`
	PartyID             uint64              `json:"partyID"`
	StartTime           string              `json:"startTime"`
	StateProvinceRegion string              `json:"stateProvinceRegion"`
	Title               string              `json:"title"`
	ZipCode             uint32              `json:"zipCode"`
}

// Host : A host of a party from the database in json format
type Host struct {
	IsMainHost bool   `json:"isMainHost"`
	Name       string `json:"name"`
}

// Invitee : An invitee to a party from the database in json format
type Invitee struct {
	IsMale        bool   `json:"isMale"`
	Name          string `json:"name"`
	Rating        string `json:"rating"`
	TimeLastRated string `json:"timeLastRated"`
	Status        string `json:"status"`
}

// BarData : A bar from the database in json format
type BarData struct {
	AddressLine1        string               `json:"addressLine1"`
	AddressLine2        string               `json:"addressLine2"`
	Attendees           map[string]*Attendee `json:"attendees"`
	BarID               uint64               `json:"barID"`
	City                string               `json:"city"`
	ClosingTime         string               `json:"closingTime"`
	Country             string               `json:"country"`
	Description         string               `json:"description"`
	Hosts               map[string]*Host     `json:"hosts"`
	LastCall            string               `json:"lastCall"`
	Latitude            float64              `json:"latitude"`
	Longitude           float64              `json:"longitude"`
	Name                string               `json:"name"`
	StateProvinceRegion string               `json:"stateProvinceRegion"`
	ZipCode             uint32               `json:"zipCode"`
}

// Attendee : An attendee to a bar in the database in json format
type Attendee struct {
	IsMale        bool   `json:"isMale"`
	Name          string `json:"name"`
	Rating        string `json:"rating"`
	Status        string `json:"status"`
	TimeLastRated string `json:"timeLastRated"`
}

// QueryResult : The result of a query in json format
type QueryResult struct {
	Succeeded bool     `json:"succeeded"`
	Errors    []string `json:"errors"`
}
