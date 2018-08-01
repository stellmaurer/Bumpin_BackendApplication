/*******************************************************
 * Copyright (C) 2018 Stephen Ellmaurer <stellmaurer@gmail.com>
 *
 * This file is part of the Bumpin mobile app project.
 *
 * The Bumpin project and any of the files within the Bumpin
 * project can not be copied and/or distributed without
 * the express permission of Stephen Ellmaurer.
 *******************************************************/

package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Change my rating for a party
func rateParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	facebookID := r.Form.Get("facebookID")
	rating := r.Form.Get("rating")
	timeLastRated := r.Form.Get("timeLastRated")
	timeOfLastKnownLocation := r.Form.Get("timeOfLastKnownLocation")
	queryResult := ratePartyHelper(partyID, facebookID, rating, timeLastRated, timeOfLastKnownLocation)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Change my rating for a bar (add my info to the attendees map if need be)
func rateBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	name := r.Form.Get("name")
	rating := r.Form.Get("rating")
	status := r.Form.Get("status")
	timeLastRated := r.Form.Get("timeLastRated")
	timeOfLastKnownLocation := r.Form.Get("timeOfLastKnownLocation")
	timeOfCheckIn := r.Form.Get("timeOfCheckIn")
	saidThereWasACover, saidThereWasACoverConvErr := strconv.ParseBool(r.Form.Get("saidThereWasACover"))
	saidThereWasALine, saidThereWasALineConvErr := strconv.ParseBool(r.Form.Get("saidThereWasALine"))
	var queryResult = QueryResult{}
	if isMaleConvErr != nil {
		queryResult.Error = "rateBar function: isMale parameter issue. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if saidThereWasACoverConvErr != nil {
		queryResult.Error = "rateBar function: HTTP post request saidThereWasACover parameter messed up. " + saidThereWasACoverConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if saidThereWasALineConvErr != nil {
		queryResult.Error = "rateBar function: HTTP post request saidThereWasALine parameter messed up. " + saidThereWasALineConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = rateBarHelper(barID, facebookID, isMale, name, rating, status, timeLastRated, timeOfLastKnownLocation, timeOfCheckIn, saidThereWasACover, saidThereWasALine)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Change my rating for a party
func clearRatingForParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	facebookID := r.Form.Get("facebookID")
	timeLastRated := r.Form.Get("timeLastRated")
	timeOfLastKnownLocation := r.Form.Get("timeOfLastKnownLocation")
	queryResult := clearRatingForPartyHelper(partyID, facebookID, timeLastRated, timeOfLastKnownLocation)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Change my rating for a bar (add my info to the attendees map if need be)
func clearRatingForBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	name := r.Form.Get("name")
	status := r.Form.Get("status")
	timeLastRated := r.Form.Get("timeLastRated")
	timeOfLastKnownLocation := r.Form.Get("timeOfLastKnownLocation")
	timeOfCheckIn := r.Form.Get("timeOfCheckIn")
	saidThereWasACover, saidThereWasACoverConvErr := strconv.ParseBool(r.Form.Get("saidThereWasACover"))
	saidThereWasALine, saidThereWasALineConvErr := strconv.ParseBool(r.Form.Get("saidThereWasALine"))
	var queryResult = QueryResult{}
	if isMaleConvErr != nil {
		queryResult.Error = "clearRatingForBar function: isMale parameter issue. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if saidThereWasACoverConvErr != nil {
		queryResult.Error = "rateBar function: HTTP post request saidThereWasACover parameter messed up. " + saidThereWasACoverConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if saidThereWasALineConvErr != nil {
		queryResult.Error = "rateBar function: HTTP post request saidThereWasALine parameter messed up. " + saidThereWasALineConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = clearRatingForBarHelper(barID, facebookID, isMale, name, status, timeLastRated, timeOfLastKnownLocation, timeOfCheckIn, saidThereWasACover, saidThereWasALine)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func ratePartyHelper(partyID string, facebookID string, rating string, timeLastRated string, timeOfLastKnownLocation string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 1)
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "ratePartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var invitees = "invitees"
	stringAtParty := "atParty"
	stringRating := "rating"
	stringTimeLastRated := "timeLastRated"
	stringTimeOfLastKnownLocation := "timeOfLastKnownLocation"
	expressionAttributeNames["#i"] = &invitees
	expressionAttributeNames["#f"] = &facebookID
	expressionAttributeNames["#a"] = &stringAtParty
	expressionAttributeNames["#r"] = &stringRating
	expressionAttributeNames["#t"] = &stringTimeLastRated
	expressionAttributeNames["#l"] = &stringTimeOfLastKnownLocation
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var atPartyAttributeValue = dynamodb.AttributeValue{}
	atPartyAttributeValue.SetBOOL(true)
	expressionValuePlaceholders[":atParty"] = &atPartyAttributeValue
	var ratingAttributeValue = dynamodb.AttributeValue{}
	ratingAttributeValue.SetS(rating)
	expressionValuePlaceholders[":rating"] = &ratingAttributeValue
	var timeLastRatedAttributeValue = dynamodb.AttributeValue{}
	timeLastRatedAttributeValue.SetS(timeLastRated)
	expressionValuePlaceholders[":timeLastRated"] = &timeLastRatedAttributeValue
	var timeOfLastKnownLocationAttributeValue = dynamodb.AttributeValue{}
	timeOfLastKnownLocationAttributeValue.SetS(timeOfLastKnownLocation)
	expressionValuePlaceholders[":timeOfLastKnownLocation"] = &timeOfLastKnownLocationAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #i.#f.#r=:rating, #i.#f.#a=:atParty, #i.#f.#t=:timeLastRated, #i.#f.#l=:timeOfLastKnownLocation"
	updateItemInput.UpdateExpression = &updateExpression
	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "ratePartyHelper function: UpdateItem error (probable cause: this party doesn't exist or you aren't invited to this party). " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func rateBarHelper(barID string, facebookID string, isMale bool, name string, rating string, status string, timeLastRated string, timeOfLastKnownLocation string, timeOfCheckIn string, saidThereWasACover bool, saidThereWasALine bool) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 1)
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "rateBarHelper function: session creation error. " + err.Error()
		return queryResult
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
	var atBarAttribute = dynamodb.AttributeValue{}
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	var timeOfLastKnownLocationAttribute = dynamodb.AttributeValue{}
	var timeOfCheckInAttribute = dynamodb.AttributeValue{}
	var saidThereWasACoverAttribute = dynamodb.AttributeValue{}
	var saidThereWasALineAttribute = dynamodb.AttributeValue{}
	atBarAttribute.SetBOOL(true)
	isMaleAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	ratingAttribute.SetS(rating)
	statusAttribute.SetS(status)
	timeLastRatedAttribute.SetS(timeLastRated)
	timeOfLastKnownLocationAttribute.SetS(timeOfLastKnownLocation)
	timeOfCheckInAttribute.SetS(timeOfCheckIn)
	saidThereWasACoverAttribute.SetBOOL(saidThereWasACover)
	saidThereWasALineAttribute.SetBOOL(saidThereWasALine)
	attendeeMap["atBar"] = &atBarAttribute
	attendeeMap["isMale"] = &isMaleAttribute
	attendeeMap["name"] = &nameAttribute
	attendeeMap["rating"] = &ratingAttribute
	attendeeMap["status"] = &statusAttribute
	attendeeMap["timeLastRated"] = &timeLastRatedAttribute
	attendeeMap["timeOfLastKnownLocation"] = &timeOfLastKnownLocationAttribute
	attendeeMap["timeOfCheckIn"] = &timeOfCheckInAttribute
	attendeeMap["saidThereWasACover"] = &saidThereWasACoverAttribute
	attendeeMap["saidThereWasALine"] = &saidThereWasALineAttribute
	attendee.SetM(attendeeMap)
	expressionValuePlaceholders[":attendee"] = &attendee

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #a.#f=:attendee"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "rateBarHelper function: UpdateItem error. (probable cause: this bar may not exist)" + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func clearRatingForPartyHelper(partyID string, facebookID string, timeLastRated string, timeOfLastKnownLocation string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 1)
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "clearRatingForPartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var invitees = "invitees"
	stringAtParty := "atParty"
	stringRating := "rating"
	stringTimeLastRated := "timeLastRated"
	stringTimeOfLastKnownLocation := "timeOfLastKnownLocation"
	expressionAttributeNames["#i"] = &invitees
	expressionAttributeNames["#f"] = &facebookID
	expressionAttributeNames["#a"] = &stringAtParty
	expressionAttributeNames["#r"] = &stringRating
	expressionAttributeNames["#t"] = &stringTimeLastRated
	expressionAttributeNames["#l"] = &stringTimeOfLastKnownLocation
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var atPartyAttributeValue = dynamodb.AttributeValue{}
	atPartyAttributeValue.SetBOOL(false)
	expressionValuePlaceholders[":atParty"] = &atPartyAttributeValue
	var ratingAttributeValue = dynamodb.AttributeValue{}
	ratingAttributeValue.SetS("None")
	expressionValuePlaceholders[":rating"] = &ratingAttributeValue
	var timeLastRatedAttributeValue = dynamodb.AttributeValue{}
	timeLastRatedAttributeValue.SetS(timeLastRated)
	expressionValuePlaceholders[":timeLastRated"] = &timeLastRatedAttributeValue
	var timeOfLastKnownLocationAttributeValue = dynamodb.AttributeValue{}
	timeOfLastKnownLocationAttributeValue.SetS(timeOfLastKnownLocation)
	expressionValuePlaceholders[":timeOfLastKnownLocation"] = &timeOfLastKnownLocationAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #i.#f.#r=:rating, #i.#f.#a=:atParty, #i.#f.#t=:timeLastRated, #i.#f.#l=:timeOfLastKnownLocation"
	updateItemInput.UpdateExpression = &updateExpression
	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "ratePartyHelper function: UpdateItem error (probable cause: this party doesn't exist or you aren't invited to this party). " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func clearRatingForBarHelper(barID string, facebookID string, isMale bool, name string, status string, timeLastRated string, timeOfLastKnownLocation string, timeOfCheckIn string, saidThereWasACover bool, saidThereWasALine bool) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 1)
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "clearRatingForBarHelper function: session creation error. " + err.Error()
		return queryResult
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
	var atBarAttribute = dynamodb.AttributeValue{}
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	var timeOfLastKnownLocationAttribute = dynamodb.AttributeValue{}
	var timeOfCheckInAttribute = dynamodb.AttributeValue{}
	var saidThereWasACoverAttribute = dynamodb.AttributeValue{}
	var saidThereWasALineAttribute = dynamodb.AttributeValue{}
	atBarAttribute.SetBOOL(false)
	isMaleAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	ratingAttribute.SetS("None")
	statusAttribute.SetS(status)
	timeLastRatedAttribute.SetS(timeLastRated)
	timeOfLastKnownLocationAttribute.SetS(timeOfLastKnownLocation)
	timeOfCheckInAttribute.SetS(timeOfCheckIn)
	saidThereWasACoverAttribute.SetBOOL(saidThereWasACover)
	saidThereWasALineAttribute.SetBOOL(saidThereWasALine)
	attendeeMap["atBar"] = &atBarAttribute
	attendeeMap["isMale"] = &isMaleAttribute
	attendeeMap["name"] = &nameAttribute
	attendeeMap["rating"] = &ratingAttribute
	attendeeMap["status"] = &statusAttribute
	attendeeMap["timeLastRated"] = &timeLastRatedAttribute
	attendeeMap["timeOfLastKnownLocation"] = &timeOfLastKnownLocationAttribute
	attendeeMap["timeOfCheckIn"] = &timeOfCheckInAttribute
	attendeeMap["saidThereWasACover"] = &saidThereWasACoverAttribute
	attendeeMap["saidThereWasALine"] = &saidThereWasALineAttribute
	attendee.SetM(attendeeMap)
	expressionValuePlaceholders[":attendee"] = &attendee

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #a.#f=:attendee"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "clearRatingForBarHelper function: UpdateItem error. (probable cause: this bar may not exist)" + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}
