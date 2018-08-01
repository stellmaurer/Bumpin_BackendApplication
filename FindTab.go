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
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func claimBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	claimKey := r.Form.Get("claimKey")
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	nameOfCreator := r.Form.Get("nameOfCreator")

	var queryResult = QueryResult{}
	if isMaleConvErr != nil {
		queryResult.Error = "claimBar function: isMale parameter issue. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}

	queryResult = claimBarHelper(claimKey, facebookID, isMale, nameOfCreator)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Find all of the parties for partyIDs passed in.
func getParties(w http.ResponseWriter, r *http.Request) {
	var queryResult = QueryResult{}
	if r.URL.Query().Get("partyIDs") == "" {
		queryResult.Succeeded = true
		queryResult.DynamodbCalls = nil
	} else {
		partyIDs := strings.Split(r.URL.Query().Get("partyIDs"), ",")
		queryResult = getPartiesHelper(partyIDs)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Find all bars that I'm close to
func barsCloseToMe(w http.ResponseWriter, r *http.Request) {
	latitude, latitudeErr := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	longitude, longitudeErr := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if latitudeErr != nil {
		queryResult.Error = "barsCloseToMe function: HTTP get request latitude parameter messed up. " + latitudeErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if longitudeErr != nil {
		queryResult.Error = "barsCloseToMe function: HTTP get request longitude parameter messed up. " + longitudeErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = barsCloseToMeHelper(latitude, longitude)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Change my attendance status to a party
func changeAttendanceStatusToParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	facebookID := r.Form.Get("facebookID")
	status := r.Form.Get("status")
	queryResult := changeAttendanceStatusToPartyHelper(partyID, facebookID, status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Change my attendance status to a bar (add my info to the attendees map if need be)
func changeAttendanceStatusToBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	facebookID := r.Form.Get("facebookID")
	atBar, atBarConvErr := strconv.ParseBool(r.Form.Get("atBar"))
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
	queryResult.Succeeded = false
	if isMaleConvErr != nil {
		queryResult.Error = "changeAttendanceStatusToBar function: HTTP post request isMale parameter messed up. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if atBarConvErr != nil {
		queryResult.Error = "changeAttendanceStatusToBar function: HTTP post request atBar parameter messed up. " + atBarConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if saidThereWasACoverConvErr != nil {
		queryResult.Error = "changeAttendanceStatusToBar function: HTTP post request saidThereWasACover parameter messed up. " + saidThereWasACoverConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if saidThereWasALineConvErr != nil {
		queryResult.Error = "changeAttendanceStatusToBar function: HTTP post request saidThereWasALine parameter messed up. " + saidThereWasALineConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = createOrUpdateAttendeeHelper(barID, facebookID, atBar, isMale, name, rating, status, timeLastRated, timeOfLastKnownLocation, timeOfCheckIn, saidThereWasACover, saidThereWasALine)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Change my atPartyStatus
func changeAtPartyStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	facebookID := r.Form.Get("facebookID")
	atParty, atPartyConvErr := strconv.ParseBool(r.Form.Get("atParty"))
	timeOfLastKnownLocation := r.Form.Get("timeOfLastKnownLocation")
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if atPartyConvErr != nil {
		queryResult.Error = "changeAtPartyStatus function: HTTP post request atParty parameter messed up. " + atPartyConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = changeAtPartyStatusHelper(partyID, facebookID, atParty, timeOfLastKnownLocation)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Change the user's atBar status
func changeAtBarStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	facebookID := r.Form.Get("facebookID")
	atBar, atBarConvErr := strconv.ParseBool(r.Form.Get("atBar"))
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
	queryResult.Succeeded = false
	if isMaleConvErr != nil {
		queryResult.Error = "changeAtBarStatus function: HTTP post request isMale parameter messed up. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if atBarConvErr != nil {
		queryResult.Error = "changeAtBarStatus function: HTTP post request atBar parameter messed up. " + atBarConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if saidThereWasACoverConvErr != nil {
		queryResult.Error = "changeAttendanceStatusToBar function: HTTP post request saidThereWasACover parameter messed up. " + saidThereWasACoverConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if saidThereWasALineConvErr != nil {
		queryResult.Error = "changeAttendanceStatusToBar function: HTTP post request saidThereWasALine parameter messed up. " + saidThereWasALineConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = createOrUpdateAttendeeHelper(barID, facebookID, atBar, isMale, name, rating, status, timeLastRated, timeOfLastKnownLocation, timeOfCheckIn, saidThereWasACover, saidThereWasALine)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Invite another friend to the party if you have invitations left.
//		A host has unlimited invitations.
func inviteFriendToParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	myFacebookID := r.Form.Get("myFacebookID")
	isHost, isHostErr := strconv.ParseBool(r.Form.Get("isHost"))
	numberOfInvitesToGive := r.Form.Get("numberOfInvitesToGive")
	friendFacebookID := r.Form.Get("friendFacebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	name := r.Form.Get("name")

	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if isMaleConvErr != nil {
		queryResult.Error = "inviteFriendToParty function: HTTP post request isMale parameter issue. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if isHostErr != nil {
		queryResult.Error = "inviteFriendToParty function: HTTP post request isHost parameter issue. " + isHostErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = inviteFriendToPartyHelper(partyID, myFacebookID, isHost, numberOfInvitesToGive, friendFacebookID, isMale, name)

	if queryResult.Succeeded == true {
		message := r.Form.Get("name") + " invited you to a party."
		sendPushNotificationsQueryResult := createAndSendNotificationsToThesePeople([]string{friendFacebookID}, message, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, sendPushNotificationsQueryResult)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func claimBarHelper(claimKey string, facebookID string, isMale bool, nameOfCreator string) QueryResult {
	barID := "-1"
	var queryResult QueryResult
	if claimKey != "AdminUe28GTttHi3L30Jjd3ILLLAdmin" {
		queryResult = getClaimKeyHelper(claimKey)
		if queryResult.Succeeded == false {
			return queryResult
		}
		barID = queryResult.Error

		queryResult = deleteBarKeyHelper(claimKey)
		if queryResult.Succeeded == false {
			return queryResult
		}
	}

	getBarQueryResult := getBar(barID)
	if getBarQueryResult.Succeeded == false {
		return getBarQueryResult
	}
	removeHostsQueryResult := removeBarFromBarHostForMapInPersonTableForEveryHostOfTheBar(barID, getBarQueryResult.Bars[0].Hosts)
	if removeHostsQueryResult.Succeeded == false {
		return removeHostsQueryResult
	}

	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 2)
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "claimBarHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)

	// Finally
	expressionAttributeNames := make(map[string]*string)
	var hostsString = "hosts"
	expressionAttributeNames["#hosts"] = &hostsString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	hostsMap := make(map[string]*dynamodb.AttributeValue)
	var hosts = dynamodb.AttributeValue{}
	hostMap := make(map[string]*dynamodb.AttributeValue)
	var host = dynamodb.AttributeValue{}
	var isMainHostAttribute = dynamodb.AttributeValue{}
	var hostStatusAttribute = dynamodb.AttributeValue{}
	var nameOfCreatorAttribute = dynamodb.AttributeValue{}
	isMainHostAttribute.SetBOOL(true)
	hostStatusAttribute.SetS("Accepted")
	nameOfCreatorAttribute.SetS(nameOfCreator)
	hostMap["isMainHost"] = &isMainHostAttribute
	hostMap["name"] = &nameOfCreatorAttribute
	hostMap["status"] = &hostStatusAttribute
	host.SetM(hostMap)
	hostsMap[facebookID] = &host
	hosts.SetM(hostsMap)
	expressionValuePlaceholders[":hosts"] = &hosts

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #hosts=:hosts"
	updateItemInput.UpdateExpression = &updateExpression
	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "claimBarHelper function: Problem adding host to Bar: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the person's information to let them
	//     know that they are hosting this bar.
	expressionAttributeNames = make(map[string]*string)
	barHostFor := "barHostFor"
	expressionAttributeNames["#barHostFor"] = &barHostFor
	expressionAttributeNames["#barID"] = &barID
	expressionValuePlaceholders = make(map[string]*dynamodb.AttributeValue)
	var barIDBoolAttribute = dynamodb.AttributeValue{}
	barIDBoolAttribute.SetBOOL(true)
	expressionValuePlaceholders[":bool"] = &barIDBoolAttribute

	keyMap = make(map[string]*dynamodb.AttributeValue)
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression = "SET #barHostFor.#barID=:bool"
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall2 = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall2.Error = "claimBarHelper function (updating person's info): UpdateItem error (probable cause: your facebookID isn't in the database). " + updateItemOutputErr.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func getPartiesHelper(partyIDs []string) QueryResult {
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
		queryResult.Error = "getPartiesHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var batchGetItemInput = dynamodb.BatchGetItemInput{}
	attributesAndValues := make([]map[string]*dynamodb.AttributeValue, len(partyIDs))
	for i := 0; i < len(partyIDs); i++ {
		var attributeValue = dynamodb.AttributeValue{}
		attributeValue.SetS(partyIDs[i])
		attributesAndValues[i] = make(map[string]*dynamodb.AttributeValue)
		attributesAndValues[i]["partyID"] = &attributeValue
	}
	var keysAndAttributes dynamodb.KeysAndAttributes
	keysAndAttributes.SetKeys(attributesAndValues)
	requestedItems := make(map[string]*dynamodb.KeysAndAttributes)
	requestedItems["Party"] = &keysAndAttributes
	batchGetItemInput.SetRequestItems(requestedItems)
	//getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"partyID": &attributeValue})
	batchGetItemOutput, err2 := getter.DynamoDB.BatchGetItem(&batchGetItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getPartiesHelper function: BatchGetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := batchGetItemOutput.Responses
	parties := make([]PartyData, len(partyIDs))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data["Party"], &parties)
	if jsonErr != nil {
		queryResult.Error = "getPartiesHelper function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.Parties = parties
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func barsCloseToMeHelper(latitude float64, longitude float64) QueryResult {
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
		queryResult.Error = "barsCloseToMeHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)

	// Finally
	var bars []BarData
	firstCall := true
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue

	for {
		var scanItemsInput = dynamodb.ScanInput{}
		scanItemsInput.SetTableName("Bar")
		if firstCall == false && lastEvaluatedKey == nil {
			break
		} else {
			scanItemsInput.SetExclusiveStartKey(lastEvaluatedKey)
		}
		expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
		// approximately (it's a square not a circle) getting bars within a 25 mile radius
		// # of degrees / 25 miles = 0.3625 degrees of latitude
		degreesOfLatitudeWhichEqual25Miles := 0.3625
		// # of degrees / 25 miles = (1 degree of longitude / (69.1703 * COS(Latitude * 0.0174533)) ) * 25 miles
		degreesOfLongitudeWhichEqual25Miles := (1 / (69.1703 * math.Cos(latitude*0.0174533))) * 25
		latitudeSouth := latitude - degreesOfLatitudeWhichEqual25Miles
		latitudeNorth := latitude + degreesOfLatitudeWhichEqual25Miles
		longitudeEast := longitude - degreesOfLongitudeWhichEqual25Miles
		longitudeWest := longitude + degreesOfLongitudeWhichEqual25Miles
		latitudeSouthAttributeValue := dynamodb.AttributeValue{}
		latitudeNorthAttributeValue := dynamodb.AttributeValue{}
		longitudeEastAttributeValue := dynamodb.AttributeValue{}
		longitudeWestAttributeValue := dynamodb.AttributeValue{}
		latitudeSouthAttributeValue.SetN(strconv.FormatFloat(latitudeSouth, 'f', -1, 64))
		latitudeNorthAttributeValue.SetN(strconv.FormatFloat(latitudeNorth, 'f', -1, 64))
		longitudeEastAttributeValue.SetN(strconv.FormatFloat(longitudeEast, 'f', -1, 64))
		longitudeWestAttributeValue.SetN(strconv.FormatFloat(longitudeWest, 'f', -1, 64))
		expressionValuePlaceholders[":latitudeSouth"] = &latitudeSouthAttributeValue
		expressionValuePlaceholders[":latitudeNorth"] = &latitudeNorthAttributeValue
		expressionValuePlaceholders[":longitudeEast"] = &longitudeEastAttributeValue
		expressionValuePlaceholders[":longitudeWest"] = &longitudeWestAttributeValue
		scanItemsInput.SetExpressionAttributeValues(expressionValuePlaceholders)
		scanItemsInput.SetFilterExpression("(latitude BETWEEN :latitudeSouth AND :latitudeNorth) AND (longitude BETWEEN :longitudeEast AND :longitudeWest)")
		scanItemsOutput, err2 := getter.DynamoDB.Scan(&scanItemsInput)
		var dynamodbCall = DynamodbCall{}
		if err2 != nil {
			dynamodbCall.Error = "barsCloseToMeHelper function: Scan error. " + err2.Error()
			dynamodbCall.Succeeded = false
			queryResult.DynamodbCalls[0] = dynamodbCall
			queryResult.Error += dynamodbCall.Error
			return queryResult
		}
		dynamodbCall.Succeeded = true
		queryResult.DynamodbCalls[0] = dynamodbCall

		data := scanItemsOutput.Items
		barsOnThisPage := make([]BarData, len(data))
		jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &barsOnThisPage)
		if jsonErr != nil {
			queryResult.Error = "barsCloseToMeHelper function: UnmarshalListOfMaps error. " + jsonErr.Error()
			return queryResult
		}
		bars = append(bars, barsOnThisPage...)
		lastEvaluatedKey = scanItemsOutput.LastEvaluatedKey
		firstCall = false
	}

	queryResult.Bars = bars
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func changeAttendanceStatusToPartyHelper(partyID string, facebookID string, status string) QueryResult {
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
		queryResult.Error = "changeAttendanceStatusToPartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var invitees = "invitees"
	stringStatus := "status"
	expressionAttributeNames["#i"] = &invitees
	expressionAttributeNames["#f"] = &facebookID
	expressionAttributeNames["#s"] = &stringStatus
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var statusAttributeValue = dynamodb.AttributeValue{}
	statusAttributeValue.SetS(status)
	expressionValuePlaceholders[":status"] = &statusAttributeValue
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #i.#f.#s=:status"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "changeAttendanceStatusToPartyHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func createOrUpdateAttendeeHelper(barID string, facebookID string, atBar bool, isMale bool, name string, rating string, status string, timeLastRated string, timeOfLastKnownLocation string, timeOfCheckIn string, saidThereWasACover bool, saidThereWasALine bool) QueryResult {
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
		queryResult.Error = "createOrUpdateAttendeeHelper function: session creation error. " + err.Error()
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
	atBarAttribute.SetBOOL(atBar)
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
		dynamodbCall.Error = "createOrUpdateAttendeeHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func changeAtPartyStatusHelper(partyID string, facebookID string, atParty bool, timeOfLastKnownLocation string) QueryResult {
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
		queryResult.Error = "changeAtPartyStatusHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var invitees = "invitees"
	stringAtParty := "atParty"
	stringTimeOfLastKnownLocation := "timeOfLastKnownLocation"
	expressionAttributeNames["#i"] = &invitees
	expressionAttributeNames["#f"] = &facebookID
	expressionAttributeNames["#a"] = &stringAtParty
	expressionAttributeNames["#t"] = &stringTimeOfLastKnownLocation
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var atPartyAttributeValue = dynamodb.AttributeValue{}
	var timeOfLastKnownLocationAttributeValue = dynamodb.AttributeValue{}
	atPartyAttributeValue.SetBOOL(atParty)
	timeOfLastKnownLocationAttributeValue.SetS(timeOfLastKnownLocation)
	expressionValuePlaceholders[":atParty"] = &atPartyAttributeValue
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
	updateExpression := "SET #i.#f.#a=:atParty, #i.#f.#t=:timeOfLastKnownLocation"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "changeAtPartyStatusHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func inviteFriendToPartyHelper(partyID string, myFacebookID string, isHost bool, numberOfInvitesToGive string, friendFacebookID string, isMale bool, name string) QueryResult {
	// invitees can't let their invitees give out their own invitations
	var numberOfInvitationsLeft = "0"
	if isHost == true {
		// hosts can let their invitees give out their own invitations
		numberOfInvitationsLeft = numberOfInvitesToGive
	}
	rating := "None"
	var status = "Invited"
	// myFacebookID will equal friendFacebookID if this function is being called
	//		by the acceptInvitationToHostParty function, so logically we want
	//		the invitation status of the new host to be "Going".
	if myFacebookID == friendFacebookID {
		status = "Going"
	}
	// constant in the past to make sure the invitee
	//     can rate the party right away
	timeLastRated := "2000-01-01T00:00:00Z"
	// doesn't really matter what time is initially put in here
	timeOfLastKnownLocation := "2000-01-01T00:00:00Z"

	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 2)
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "inviteFriendToPartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var invitees = "invitees"
	numberOfInvitationsLeftString := "numberOfInvitationsLeft"
	expressionAttributeNames["#i"] = &invitees
	expressionAttributeNames["#f"] = &friendFacebookID
	if isHost == false {
		// These are only relevant if a non-host is inviting this person
		expressionAttributeNames["#m"] = &myFacebookID
		expressionAttributeNames["#n"] = &numberOfInvitationsLeftString
	}
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)

	inviteeMap := make(map[string]*dynamodb.AttributeValue)
	var invitee = dynamodb.AttributeValue{}
	var atPartyAttribute = dynamodb.AttributeValue{}
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var numberOfInvitationsLeftAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	var timeOfLastKnownLocationAttribute = dynamodb.AttributeValue{}
	atPartyAttribute.SetBOOL(false)
	isMaleAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	numberOfInvitationsLeftAttribute.SetN(numberOfInvitationsLeft)
	ratingAttribute.SetS(rating)
	statusAttribute.SetS(status)
	timeLastRatedAttribute.SetS(timeLastRated)
	timeOfLastKnownLocationAttribute.SetS(timeOfLastKnownLocation)
	inviteeMap["atParty"] = &atPartyAttribute
	inviteeMap["isMale"] = &isMaleAttribute
	inviteeMap["name"] = &nameAttribute
	inviteeMap["numberOfInvitationsLeft"] = &numberOfInvitationsLeftAttribute
	inviteeMap["rating"] = &ratingAttribute
	inviteeMap["status"] = &statusAttribute
	inviteeMap["timeLastRated"] = &timeLastRatedAttribute
	inviteeMap["timeOfLastKnownLocation"] = &timeOfLastKnownLocationAttribute
	invitee.SetM(inviteeMap)
	expressionValuePlaceholders[":invitee"] = &invitee

	var decrementAttribute = dynamodb.AttributeValue{}
	decrementAttribute.SetN("-1")
	var oneAttribute = dynamodb.AttributeValue{}
	oneAttribute.SetN("1")
	if isHost == false {
		// These are only relevant if a non-host is inviting this person
		expressionValuePlaceholders[":decrement"] = &decrementAttribute
		expressionValuePlaceholders[":one"] = &oneAttribute
	}

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	var updateExpression = ""
	if isHost == true {
		// If this is a host inviting someone, ignore the number of invites they
		//		have left, because they get unlimited invitations as a host.
		updateExpression = "SET #i.#f=:invitee"
	} else {
		var conditionExpression = "#i.#m.#n >= :one"
		updateItemInput.SetConditionExpression(conditionExpression)
		updateExpression = "SET #i.#f=:invitee ADD #i.#m.#n :decrement"
	}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr1 := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall1 = DynamodbCall{}
	if updateItemOutputErr1 != nil {
		dynamodbCall1.Error = "inviteFriendToPartyHelper function: UpdateItem1 error (probable cause: you either don't have any invitations left or your friend is already invited to this party). " + updateItemOutputErr1.Error()
		dynamodbCall1.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall1
		return queryResult
	}
	dynamodbCall1.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall1
	// Now we need to update the friend's information to let them
	//     know that they are invited to this party.
	expressionAttributeNames2 := make(map[string]*string)
	invitedTo := "invitedTo"
	expressionAttributeNames2["#i"] = &invitedTo
	expressionAttributeNames2["#p"] = &partyID
	expressionValuePlaceholders2 := make(map[string]*dynamodb.AttributeValue)
	var friendFacebookIDBoolAttribute = dynamodb.AttributeValue{}
	friendFacebookIDBoolAttribute.SetBOOL(true)
	expressionValuePlaceholders2[":bool"] = &friendFacebookIDBoolAttribute

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(friendFacebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetExpressionAttributeValues(expressionValuePlaceholders2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	updateExpression2 := "SET #i.#p=:bool"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, updateItemOutputErr2 := getter.DynamoDB.UpdateItem(&updateItemInput2)

	var dynamodbCall2 = DynamodbCall{}
	if updateItemOutputErr2 != nil {
		dynamodbCall2.Error = "inviteFriendToPartyHelper function: UpdateItem2 error (probable cause: your friend's facebookID isn't in the database). " + updateItemOutputErr2.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func removeFriendFromPartyHelper(partyID string, friendFacebookID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 2)
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "removeFriendFromPartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var invitees = "invitees"
	expressionAttributeNames["#i"] = &invitees
	expressionAttributeNames["#f"] = &friendFacebookID

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	var updateExpression = "REMOVE #i.#f"
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr1 := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall1 = DynamodbCall{}
	if updateItemOutputErr1 != nil {
		dynamodbCall1.Error = "removeFriendFromPartyHelper function: UpdateItem1 error. " + updateItemOutputErr1.Error()
		dynamodbCall1.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall1
		return queryResult
	}
	dynamodbCall1.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall1
	// Now we need to update the friend's information to let them
	//     know that they aren't invited to this party.
	expressionAttributeNames2 := make(map[string]*string)
	var invitedTo = "invitedTo"
	expressionAttributeNames2["#i"] = &invitedTo
	expressionAttributeNames2["#p"] = &partyID

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(friendFacebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	updateExpression2 := "REMOVE #i.#p"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, updateItemOutputErr2 := getter.DynamoDB.UpdateItem(&updateItemInput2)

	var dynamodbCall2 = DynamodbCall{}
	if updateItemOutputErr2 != nil {
		dynamodbCall2.Error = "removeFriendFromPartyHelper function: UpdateItem2 error (probable cause: your friend's facebookID isn't in the database). " + updateItemOutputErr2.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func sendInvitationsAsGuestOfParty(w http.ResponseWriter, r *http.Request) {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	guestFacebookID := r.Form.Get("guestFacebookID")
	var additionsListFacebookID []string
	var additionsListIsMaleString []string
	var additionsListName []string

	if r.Form.Get("additionsListFacebookID") != "" {
		additionsListFacebookID = strings.Split(r.Form.Get("additionsListFacebookID"), ",")
	}
	if r.Form.Get("additionsListIsMale") != "" {
		additionsListIsMaleString = strings.Split(r.Form.Get("additionsListIsMale"), ",")
	}
	// convert IsMale string array to IsMale bool array
	var additionsListIsMale = make([]bool, len(additionsListIsMaleString))
	for i := 0; i < len(additionsListIsMaleString); i++ {
		isMale, isMaleConvErr := strconv.ParseBool(additionsListIsMaleString[i])
		if isMaleConvErr != nil {
			queryResult.Error = "sendInvitationsAsGuestOfParty function: HTTP post request isMale parameter issue. " + isMaleConvErr.Error()
			json.NewEncoder(w).Encode(queryResult)
			return
		}
		additionsListIsMale[i] = isMale
	}
	if r.Form.Get("additionsListName") != "" {
		additionsListName = strings.Split(r.Form.Get("additionsListName"), ",")
	}
	if (len(additionsListFacebookID) != len(additionsListIsMale)) || (len(additionsListIsMale) != len(additionsListName)) {
		queryResult.Error = "sendInvitationsAsGuestOfParty function: HTTP post request parameter issues (additions lists): facebookID array, isMale array, and name array aren't the same length."
		json.NewEncoder(w).Encode(queryResult)
		return
	}

	queryResult = sendInvitationsAsGuestOfPartyHelper(partyID, guestFacebookID, additionsListFacebookID, additionsListIsMale, additionsListName)

	if queryResult.Succeeded == true {
		message := r.Form.Get("guestName") + " invited you to a party."
		sendPushNotificationsQueryResult := createAndSendNotificationsToThesePeople(additionsListFacebookID, message, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, sendPushNotificationsQueryResult)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func sendInvitationsAsGuestOfPartyHelper(partyID string, guestFacebookID string,
	additionsListFacebookID []string, additionsListIsMale []bool, additionsListName []string) QueryResult {
	var numberOfInvitesToGive = "0"
	var isHost = false

	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	var inviteFriendQueryResult = QueryResult{}
	for i := 0; i < len(additionsListFacebookID); i++ {
		inviteFriendQueryResult = inviteFriendToPartyHelper(partyID, guestFacebookID, isHost, numberOfInvitesToGive, additionsListFacebookID[i], additionsListIsMale[i], additionsListName[i])
		if inviteFriendQueryResult.Succeeded == false {
			queryResult = convertTwoQueryResultsToOne(inviteFriendQueryResult, queryResult)
		}
	}
	return queryResult
}
