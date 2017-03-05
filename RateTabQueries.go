package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Change my rating for a party
func rateParty(w http.ResponseWriter, r *http.Request) {
	partyID := r.URL.Query().Get("partyID")
	facebookID := r.URL.Query().Get("facebookID")
	rating := r.URL.Query().Get("rating")
	timeLastRated := r.URL.Query().Get("timeLastRated")

	data, err := ratePartyHelper(partyID, facebookID, rating, timeLastRated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var party = PartyData{}
	jsonErr := dynamodbattribute.UnmarshalMap(data, &party)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(party)
}

func ratePartyHelper(partyID string, facebookID string, rating string, timeLastRated string) (map[string]*dynamodb.AttributeValue, error) {
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
	expressionAttributeNames := make(map[string]*string)
	var invitees = "invitees"
	stringRating := "rating"
	stringTimeLastRated := "timeLastRated"
	expressionAttributeNames["#i"] = &invitees
	expressionAttributeNames["#f"] = &facebookID
	expressionAttributeNames["#r"] = &stringRating
	expressionAttributeNames["#t"] = &stringTimeLastRated
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var ratingAttributeValue = dynamodb.AttributeValue{}
	ratingAttributeValue.SetS(rating)
	expressionValuePlaceholders[":rating"] = &ratingAttributeValue
	var timeLastRatedAttributeValue = dynamodb.AttributeValue{}
	timeLastRatedAttributeValue.SetS(timeLastRated)
	expressionValuePlaceholders[":timeLastRated"] = &timeLastRatedAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #i.#f.#r=:rating, #i.#f.#t=:timeLastRated"
	updateItemInput.UpdateExpression = &updateExpression

	updateItemOutput, err := getter.DynamoDB.UpdateItem(&updateItemInput)
	return updateItemOutput.Attributes, err
}

// Change my rating for a bar (add my info to the attendees map if need be)
func rateBar(w http.ResponseWriter, r *http.Request) {
	barID := r.URL.Query().Get("barID")
	facebookID := r.URL.Query().Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.URL.Query().Get("isMale"))
	name := r.URL.Query().Get("name")
	rating := r.URL.Query().Get("rating")
	status := "T"
	timeLastRated := r.URL.Query().Get("timeLastRated")

	if isMaleConvErr != nil {
		http.Error(w, "Parameter issue: "+isMaleConvErr.Error(), http.StatusInternalServerError)
		return
	}

	data, err := rateBarHelper(barID, facebookID, isMale, name, rating, status, timeLastRated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var bar = BarData{}
	jsonErr := dynamodbattribute.UnmarshalMap(data, &bar)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(bar)
}

func rateBarHelper(barID string, facebookID string, isMale bool, name string, rating string, status string, timeLastRated string) (map[string]*dynamodb.AttributeValue, error) {
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
	expressionAttributeNames := make(map[string]*string)
	attendees := "attendees"
	expressionAttributeNames["#a"] = &attendees
	expressionAttributeNames["#f"] = &facebookID
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)

	attendeeMap := make(map[string]*dynamodb.AttributeValue)
	var attendee = dynamodb.AttributeValue{}
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	isMaleAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	ratingAttribute.SetS(rating)
	statusAttribute.SetS(status)
	timeLastRatedAttribute.SetS(timeLastRated)
	attendeeMap["isMale"] = &isMaleAttribute
	attendeeMap["name"] = &nameAttribute
	attendeeMap["rating"] = &ratingAttribute
	attendeeMap["status"] = &statusAttribute
	attendeeMap["timeLastRated"] = &timeLastRatedAttribute
	attendee.SetM(attendeeMap)
	expressionValuePlaceholders[":attendee"] = &attendee

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #a.#f=:attendee"
	updateItemInput.UpdateExpression = &updateExpression

	updateItemOutput, err := getter.DynamoDB.UpdateItem(&updateItemInput)
	return updateItemOutput.Attributes, err
}
