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
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func createOrUpdatePerson(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	name := r.Form.Get("name")
	platform := r.Form.Get("platform")
	deviceToken := r.Form.Get("deviceToken")
	sevenPMLocalHourInZulu := r.Form.Get("sevenPMLocalHourInZulu")
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if isMaleConvErr != nil {
		queryResult.Error = "createOrUpdatePerson function: HTTP get request isMale parameter messed up. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult1 := createPersonHelper(facebookID, isMale, name, platform, deviceToken, sevenPMLocalHourInZulu)
	queryResult2 := updatePersonHelper(facebookID, isMale, name, platform, deviceToken, sevenPMLocalHourInZulu)
	queryResult = convertTwoQueryResultsToOne(queryResult1, queryResult2)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	facebookID := r.URL.Query().Get("facebookID")
	queryResult := getPersonHelper(facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func createPersonHelper(facebookID string, isMale bool, name string, platform string, deviceToken string, sevenPMLocalHourInZulu string) QueryResult {
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
		queryResult.Error = "createPersonHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var facebookIDString = "facebookID"
	expressionAttributeNames["#fbid"] = &facebookIDString

	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var barHostForAttributeValue = dynamodb.AttributeValue{}
	var facebookIDAttributeValue = dynamodb.AttributeValue{}
	var invitedToAttributeValue = dynamodb.AttributeValue{}
	var isMaleAttributeValue = dynamodb.AttributeValue{}
	var nameAttributeValue = dynamodb.AttributeValue{}
	var outstandingNotificationsAttributeValue = dynamodb.AttributeValue{}
	var platformAttributeValue = dynamodb.AttributeValue{}
	var deviceTokenAttributeValue = dynamodb.AttributeValue{}
	var partyHostForAttributeValue = dynamodb.AttributeValue{}
	var peopleBlockingTheirActivityFromMeAttributeValue = dynamodb.AttributeValue{}
	var peopleToIgnoreAttributeValue = dynamodb.AttributeValue{}
	var sevenPMLocalHourInZuluAttributeValue = dynamodb.AttributeValue{}
	var numberOfFriendsThatMightGoOutAttributeValue = dynamodb.AttributeValue{}
	var goingOutStatusAttributeValue = dynamodb.AttributeValue{}
	barHostForAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	facebookIDAttributeValue.SetS(facebookID)
	invitedToAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	isMaleAttributeValue.SetBOOL(isMale)
	nameAttributeValue.SetS(name)
	outstandingNotificationsAttributeValue.SetN("0")
	platformAttributeValue.SetS(platform)
	deviceTokenAttributeValue.SetS(deviceToken)
	partyHostForAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	peopleBlockingTheirActivityFromMeAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	peopleToIgnoreAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	sevenPMLocalHourInZuluAttributeValue.SetN(sevenPMLocalHourInZulu)
	numberOfFriendsThatMightGoOutAttributeValue.SetN("0")
	statusMap := make(map[string]*dynamodb.AttributeValue)
	var manuallySetAttribute = dynamodb.AttributeValue{}
	var goingOutAttribute = dynamodb.AttributeValue{}
	var timeGoingOutStatusWasSetAttribute = dynamodb.AttributeValue{}
	manuallySetAttribute.SetS("No")
	goingOutAttribute.SetS("Unknown")
	timeGoingOutStatusWasSetAttribute.SetS("2000-01-01T00:00:00Z")
	statusMap["manuallySet"] = &manuallySetAttribute
	statusMap["goingOut"] = &goingOutAttribute
	statusMap["timeGoingOutStatusWasSet"] = &timeGoingOutStatusWasSetAttribute
	goingOutStatusAttributeValue.SetM(statusMap)
	expressionValues["barHostFor"] = &barHostForAttributeValue
	expressionValues["facebookID"] = &facebookIDAttributeValue
	expressionValues["invitedTo"] = &invitedToAttributeValue
	expressionValues["isMale"] = &isMaleAttributeValue
	expressionValues["name"] = &nameAttributeValue
	expressionValues["outstandingNotifications"] = &outstandingNotificationsAttributeValue
	expressionValues["platform"] = &platformAttributeValue
	expressionValues["deviceToken"] = &deviceTokenAttributeValue
	expressionValues["partyHostFor"] = &partyHostForAttributeValue
	expressionValues["peopleBlockingTheirActivityFromMe"] = &peopleBlockingTheirActivityFromMeAttributeValue
	expressionValues["peopleToIgnore"] = &peopleToIgnoreAttributeValue
	expressionValues["sevenPMLocalHourInZulu"] = &sevenPMLocalHourInZuluAttributeValue
	expressionValues["numberOfFriendsThatMightGoOut"] = &numberOfFriendsThatMightGoOutAttributeValue
	expressionValues["status"] = &goingOutStatusAttributeValue

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetConditionExpression("attribute_not_exists(#fbid)")
	putItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	putItemInput.SetTableName("Person")
	putItemInput.SetItem(expressionValues)

	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "createPersonHelper function: PutItem error (this error should be seen if the person is already in the database. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
	} else {
		queryResult.DynamodbCalls = nil
	}
	queryResult.Succeeded = true
	return queryResult
}

func updatePersonHelper(facebookID string, isMale bool, name string, platform string, deviceToken string, sevenPMLocalHourInZulu string) QueryResult {
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
		queryResult.Error = "updatePersonHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var isMaleString = "isMale"
	var nameString = "name"
	var platformString = "platform"
	var deviceTokenString = "deviceToken"
	var sevenPMLocalHourInZuluString = "sevenPMLocalHourInZulu"
	expressionAttributeNames["#isMale"] = &isMaleString
	expressionAttributeNames["#name"] = &nameString
	expressionAttributeNames["#platform"] = &platformString
	expressionAttributeNames["#deviceToken"] = &deviceTokenString
	expressionAttributeNames["#sevenPMLocalHourInZulu"] = &sevenPMLocalHourInZuluString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var isMaleAttributeValue = dynamodb.AttributeValue{}
	var nameAttributeValue = dynamodb.AttributeValue{}
	var platformAttributeValue = dynamodb.AttributeValue{}
	var deviceTokenAttributeValue = dynamodb.AttributeValue{}
	var sevenPMLocalHourInZuluAttributeValue = dynamodb.AttributeValue{}
	isMaleAttributeValue.SetBOOL(isMale)
	nameAttributeValue.SetS(name)
	platformAttributeValue.SetS(platform)
	deviceTokenAttributeValue.SetS(deviceToken)
	sevenPMLocalHourInZuluAttributeValue.SetN(sevenPMLocalHourInZulu)
	expressionValuePlaceholders[":isMale"] = &isMaleAttributeValue
	expressionValuePlaceholders[":name"] = &nameAttributeValue
	expressionValuePlaceholders[":platform"] = &platformAttributeValue
	expressionValuePlaceholders[":deviceToken"] = &deviceTokenAttributeValue
	expressionValuePlaceholders[":sevenPMLocalHourInZulu"] = &sevenPMLocalHourInZuluAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "SET #isMale=:isMale, #name=:name, #platform=:platform, #deviceToken=:deviceToken, #sevenPMLocalHourInZulu=:sevenPMLocalHourInZulu"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "updatePersonHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func getPersonHelper(facebookID string) QueryResult {
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
		queryResult.Error = "getPersonHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	getItemInput.SetTableName("Person")
	var attributeValue = dynamodb.AttributeValue{}
	attributeValue.SetS(facebookID)
	getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"facebookID": &attributeValue})
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getPersonHelper function: GetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := getItemOutput.Item
	var person PersonData
	jsonErr := dynamodbattribute.UnmarshalMap(data, &person)
	if jsonErr != nil {
		queryResult.Error = "getPersonHelper function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.People = make([]PersonData, 1)
	queryResult.People[0] = person
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}
