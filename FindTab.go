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

// Find all parties I'm invited to
func myParties(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.URL.Query().Get("partyID"), ",")
	queryResult := myPartiesHelper(params)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	json.NewEncoder(w).Encode(queryResult)
}

// Change my attendance status to a bar (add my info to the attendees map if need be)
func changeAttendanceStatusToBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	name := r.Form.Get("name")
	rating := r.Form.Get("rating")
	status := r.Form.Get("status")
	timeLastRated := r.Form.Get("timeLastRated")

	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if isMaleConvErr != nil {
		queryResult.Error = "changeAttendanceStatusToBar function: HTTP post request isMale parameter messed up. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = changeAttendanceStatusToBarHelper(barID, facebookID, isMale, name, rating, status, timeLastRated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func myPartiesHelper(params []string) QueryResult {
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
		queryResult.Error = "myPartiesHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var batchGetItemInput = dynamodb.BatchGetItemInput{}
	attributesAndValues := make([]map[string]*dynamodb.AttributeValue, len(params))
	for i := 0; i < len(params); i++ {
		var attributeValue = dynamodb.AttributeValue{}
		attributeValue.SetN(params[i])
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
		dynamodbCall.Error = "myPartiesHelper function: BatchGetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := batchGetItemOutput.Responses
	parties := make([]PartyData, len(params))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data["Party"], &parties)
	if jsonErr != nil {
		queryResult.Error = "myPartiesHelper function: UnmarshalListOfMaps error. " + jsonErr.Error()
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
		queryResult.Error = "findBarsCloseToMe function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var scanItemsInput = dynamodb.ScanInput{}
	scanItemsInput.SetTableName("Bar")
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	// approximately getting bars within a 100 mile radius
	// # of degrees / 100 miles = 1.45 degrees of latitude
	degreesOfLatitudeWhichEqual100Miles := 1.45
	// # of degrees / 100 miles = (1 degree of longitude / (69.1703 * COS(Latitude * 0.0174533)) ) * 100 miles
	degreesOfLongitudeWhichEqual100Miles := (1 / (69.1703 * math.Cos(latitude*0.0174533))) * 100
	latitudeSouth := latitude - degreesOfLatitudeWhichEqual100Miles
	latitudeNorth := latitude + degreesOfLatitudeWhichEqual100Miles
	longitudeEast := longitude - degreesOfLongitudeWhichEqual100Miles
	longitudeWest := longitude + degreesOfLongitudeWhichEqual100Miles
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
		dynamodbCall.Error = "findBarsCloseToMe function: Scan error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := scanItemsOutput.Items
	bars := make([]BarData, len(data))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &bars)
	if jsonErr != nil {
		queryResult.Error = "findBarsCloseToMe function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
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
	key.SetN(partyID)
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
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func changeAttendanceStatusToBarHelper(barID string, facebookID string, isMale bool, name string, rating string, status string, timeLastRated string) QueryResult {
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
		queryResult.Error = "changeAttendanceStatusToBarHelper function: session creation error. " + err.Error()
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

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "changeAttendanceStatusToBarHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
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
	rating := "none"
	var status = "invited"
	// myFacebookID will equal friendFacebookID if this function is being called
	//		by the acceptInvitationToHostParty function, so logically we want
	//		the invitation status of the new host to be "going".
	if myFacebookID == friendFacebookID {
		status = "going"
	}
	// constant in the past to make sure the invitee
	//     can rate the party right away
	timeLastRated := "01/01/2000 00:00:00"

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
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var numberOfInvitationsLeftAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	isMaleAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	numberOfInvitationsLeftAttribute.SetN(numberOfInvitationsLeft)
	ratingAttribute.SetS(rating)
	statusAttribute.SetS(status)
	timeLastRatedAttribute.SetS(timeLastRated)
	inviteeMap["isMale"] = &isMaleAttribute
	inviteeMap["name"] = &nameAttribute
	inviteeMap["numberOfInvitationsLeft"] = &numberOfInvitationsLeftAttribute
	inviteeMap["rating"] = &ratingAttribute
	inviteeMap["status"] = &statusAttribute
	inviteeMap["timeLastRated"] = &timeLastRatedAttribute
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
	key.SetN(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	var conditionExpression = "attribute_not_exists(#i.#f) AND (#i.#m.#n >= :one)"
	var updateExpression = "SET #i.#f=:invitee ADD #i.#m.#n :decrement"
	if isHost == true {
		// If this is a host inviting someone, ignore the number of invites they
		//		have left, because they get unlimited invitations as a host.
		conditionExpression = "attribute_not_exists(#i.#f)"
		updateExpression = "SET #i.#f=:invitee"
	}
	updateItemInput.SetConditionExpression(conditionExpression)
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
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}
