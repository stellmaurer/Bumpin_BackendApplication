package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func tables(w http.ResponseWriter, r *http.Request) {
	data, err := getTables()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
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
	AddressLine1          string              `json:"addressLine1"`
	AddressLine2          string              `json:"addressLine2"`
	City                  string              `json:"city"`
	Country               string              `json:"country"`
	Details               string              `json:"details"`
	DrinksProvided        bool                `json:"drinksProvided"`
	EndTime               string              `json:"endTime"`
	FeeForDrinks          bool                `json:"feeForDrinks"`
	Hosts                 map[string]*Host    `json:"hosts"`
	Invitees              map[string]*Invitee `json:"invitees"`
	Latitude              float64             `json:"latitude"`
	Longitude             float64             `json:"longitude"`
	InvitesForNewInvitees uint16              `json:"invitesForNewInvitees"`
	PartyID               uint64              `json:"partyID"`
	StartTime             string              `json:"startTime"`
	StateProvinceRegion   string              `json:"stateProvinceRegion"`
	Title                 string              `json:"title"`
	ZipCode               uint32              `json:"zipCode"`
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
	Details             string               `json:"details"`
	Hosts               map[string]*Host     `json:"hosts"`
	LastCall            string               `json:"lastCall"`
	Latitude            float64              `json:"latitude"`
	Longitude           float64              `json:"longitude"`
	Name                string               `json:"name"`
	OpenAt              string               `json:"openAt"`
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
	Succeeded     bool           `json:"succeeded"`
	Error         string         `json:"error"`
	DynamodbCalls []DynamodbCall `json:"dynamodbCalls"`
	People        []PersonData   `json:"people"`
	Parties       []PartyData    `json:"parties"`
	Bars          []BarData      `json:"bars"`
}

// DynamodbCall : The result of a dynamodb call.
type DynamodbCall struct {
	Succeeded bool   `json:"succeeded"`
	Error     string `json:"error"`
}

func convertTwoQueryResultsToOne(queryResult1 QueryResult, queryResult2 QueryResult) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = queryResult1.Succeeded && queryResult2.Succeeded
	queryResult.Error = queryResult1.Error + " " + queryResult2.Error
	for i := 0; i < len(queryResult1.DynamodbCalls); i++ {
		queryResult.DynamodbCalls = append(queryResult.DynamodbCalls, queryResult1.DynamodbCalls[i])
	}
	for i := 0; i < len(queryResult2.DynamodbCalls); i++ {
		queryResult.DynamodbCalls = append(queryResult.DynamodbCalls, queryResult2.DynamodbCalls[i])
	}
	for i := 0; i < len(queryResult1.People); i++ {
		queryResult.People = append(queryResult.People, queryResult1.People[i])
	}
	for i := 0; i < len(queryResult2.People); i++ {
		queryResult.People = append(queryResult.People, queryResult2.People[i])
	}

	for i := 0; i < len(queryResult1.Parties); i++ {
		queryResult.Parties = append(queryResult.Parties, queryResult1.Parties[i])
	}
	for i := 0; i < len(queryResult2.Parties); i++ {
		queryResult.Parties = append(queryResult.Parties, queryResult2.Parties[i])
	}

	for i := 0; i < len(queryResult1.Bars); i++ {
		queryResult.Bars = append(queryResult.Bars, queryResult1.Bars[i])
	}
	for i := 0; i < len(queryResult2.Parties); i++ {
		queryResult.Bars = append(queryResult.Bars, queryResult2.Bars[i])
	}
	return queryResult
}

func getRandomID() uint64 {
	return uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
}
