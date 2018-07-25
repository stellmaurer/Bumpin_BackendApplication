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
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NaySoftware/go-fcm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

const (
	serverKey = "AAAAo_YT2fc:APA91bEV1ctVnAhvWzO7uOpuMBcHpwYu1LaGDgHF3KZ4GtdY1yocH90Vc_fvFlmtGDKib1vYA24ci5QUdaoozpeI_kfd9QdHwGS2L8JNDd6AZh1I-zGZ8COLEPp75c_wlAG_iFE1NbIZ"
)

func sendGoingOutStatusNotificationToPeopleWhoHaveFriendsGoingOutAndHaveALocalTimeEqualToSevenPM(w http.ResponseWriter, r *http.Request) {
	queryResult, people := findAllPeopleWhereTheirLocalTimeIsSevenPM()

	for i := 0; i < len(people); i++ {
		if people[i].NumberOfFriendsThatMightGoOut > 0 {
			var message = strconv.FormatUint(people[i].NumberOfFriendsThatMightGoOut, 10) + " of your friends might go out tonight."
			var createAndSendNotificationToThisPersonQueryResult = createAndSendNotificationToThisPerson(people[i].FacebookID, message, "-1")
			if createAndSendNotificationToThisPersonQueryResult.Succeeded == false {
				queryResult = convertTwoQueryResultsToOne(queryResult, createAndSendNotificationToThisPersonQueryResult)
			}
		}
	}

	if queryResult.Succeeded == true {
		queryResult.DynamodbCalls = nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func findAllPeopleWhereTheirLocalTimeIsSevenPM() (QueryResult, []TinyPerson) {
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
		queryResult.Error = "findAllPeopleWhereTheirLocalTimeIsSevenPM function: session creation error. " + err.Error()
		return queryResult, nil
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)

	currentTimeInZulu := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	zuluHourOfCurrentTimeString := currentTimeInZulu[11:13]

	var people []TinyPerson
	firstCall := true
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue

	for {
		var scanItemsInput = dynamodb.ScanInput{}
		scanItemsInput.SetTableName("Person")
		if firstCall == false && lastEvaluatedKey == nil {
			break
		} else {
			scanItemsInput.SetExclusiveStartKey(lastEvaluatedKey)
		}

		expressionAttributeNames := make(map[string]*string)
		var sevenPMLocalHourInZuluString = "sevenPMLocalHourInZulu"
		expressionAttributeNames["#sevenPMLocalHourInZulu"] = &sevenPMLocalHourInZuluString

		expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
		zuluHourOfCurrentTimeAttributeValue := dynamodb.AttributeValue{}
		zuluHourOfCurrentTimeAttributeValue.SetN(zuluHourOfCurrentTimeString)
		expressionValuePlaceholders[":zuluHourOfCurrentTime"] = &zuluHourOfCurrentTimeAttributeValue
		scanItemsInput.SetTableName("Person")
		scanItemsInput.SetExpressionAttributeNames(expressionAttributeNames)
		scanItemsInput.SetExpressionAttributeValues(expressionValuePlaceholders)
		scanItemsInput.SetFilterExpression("#sevenPMLocalHourInZulu = :zuluHourOfCurrentTime")
		scanItemsOutput, err2 := getter.DynamoDB.Scan(&scanItemsInput)

		var dynamodbCall = DynamodbCall{}
		if err2 != nil {
			dynamodbCall.Error = "findAllPeopleWhereTheirLocalTimeIsSevenPM function: Scan error. " + err2.Error()
			dynamodbCall.Succeeded = false
			queryResult.DynamodbCalls[0] = dynamodbCall
			queryResult.Error += dynamodbCall.Error
			return queryResult, nil
		}
		dynamodbCall.Succeeded = true
		queryResult.DynamodbCalls[0] = dynamodbCall

		data := scanItemsOutput.Items
		peopleOnThisPage := make([]TinyPerson, len(data))
		jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &peopleOnThisPage)
		if jsonErr != nil {
			queryResult.Error = "findAllPeopleWhereTheirLocalTimeIsSevenPM function: UnmarshalListOfMaps error. " + jsonErr.Error()
			return queryResult, nil
		}
		people = append(people, peopleOnThisPage...)
		lastEvaluatedKey = scanItemsOutput.LastEvaluatedKey
		firstCall = false
	}

	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult, people
}

// TinyPerson : less info than a normal person object to conserve resources
type TinyPerson struct {
	FacebookID                    string `json:"facebookID"`
	NumberOfFriendsThatMightGoOut uint64 `json:"numberOfFriendsThatMightGoOut"`
}

func clearOutstandingNotificationCountForPerson(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	facebookID := r.Form.Get("facebookID")
	queryResult := clearOutstandingNotificationCountForPersonHelper(facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func clearOutstandingNotificationCountForPersonHelper(facebookID string) QueryResult {
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
		queryResult.Error = "clearOutstandingNotificationCountForPersonHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var outstandingNotificationsString = "outstandingNotifications"
	expressionAttributeNames["#outstandingNotifications"] = &outstandingNotificationsString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var zeroAttribute = dynamodb.AttributeValue{}
	zeroAttribute.SetN("0")
	expressionValuePlaceholders[":zero"] = &zeroAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	var updateExpression = "SET #outstandingNotifications=:zero"
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall1 = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall1.Error = "clearOutstandingNotificationCountForPersonHelper function: " + updateItemOutputErr.Error()
		dynamodbCall1.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall1
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func incrementOutstandingNotificationCountForPerson(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	facebookID := r.Form.Get("facebookID")
	queryResult := incrementOutstandingNotificationCountForPersonHelper(facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func incrementOutstandingNotificationCountForPersonHelper(facebookID string) QueryResult {
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
		queryResult.Error = "incrementOutstandingNotificationCountForPersonHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var outstandingNotificationsString = "outstandingNotifications"
	expressionAttributeNames["#outstandingNotifications"] = &outstandingNotificationsString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var incrementAttribute = dynamodb.AttributeValue{}
	incrementAttribute.SetN("1")
	expressionValuePlaceholders[":increment"] = &incrementAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	var updateExpression = "ADD #outstandingNotifications :increment"
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall1 = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall1.Error = "incrementOutstandingNotificationCountForPersonHelper function: " + updateItemOutputErr.Error()
		dynamodbCall1.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall1
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func deleteNotification(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	notificationID := r.Form.Get("notificationID")
	queryResult := deleteNotificationHelper(notificationID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func deleteNotificationHelper(notificationID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false

	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "deleteNotification function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(notificationID)
	keyMap["notificationID"] = &key

	var deleteItemInput = dynamodb.DeleteItemInput{}
	deleteItemInput.SetTableName("Notification")
	deleteItemInput.SetKey(keyMap)

	_, err2 := getter.DynamoDB.DeleteItem(&deleteItemInput)
	if err2 != nil {
		queryResult.Error = "deleteNotification function: DeleteItem error. " + err2.Error()
		return queryResult
	}
	queryResult.Succeeded = true
	return queryResult
}

func testCreateAndSendNotificationsToThesePeople(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	people := []string{"10155227369101712", "10216576646672295", "10154326505409816"}
	queryResult := createAndSendNotificationsToThesePeople(people, "This is testing the createAndSendNotifications function.", "9999")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func createAndSendNotificationToThisPerson(facebookID string, message string, partyOrBarID string) QueryResult {
	var queryResult = QueryResult{}
	dynamodbCalls := make([]DynamodbCall, 0)

	var intermediateQueryResult QueryResult
	// update person's number of outstanding push notifications before we send
	//		the push notification payload so that we send the correct badger number
	intermediateQueryResult = QueryResult{}
	intermediateQueryResult.Succeeded = true
	intermediateQueryResult = incrementOutstandingNotificationCountForPersonHelper(facebookID)
	if intermediateQueryResult.Succeeded == false {
		dynamodbCalls = append(dynamodbCalls, DynamodbCall{Succeeded: false, Error: intermediateQueryResult.Error})
	}

	queryResult = getPersonHelper(facebookID)
	if queryResult.Succeeded == false {
		return queryResult
	}
	person := queryResult.People[0]
	intermediateQueryResult = QueryResult{}
	intermediateQueryResult.Succeeded = true
	// (Step 1) Send Push Notification
	if person.Platform != "Unknown" && person.DeviceToken != "Unknown" {
		if person.Platform == "iOS" {
			intermediateQueryResult = sendiOSPushNotification(person.DeviceToken, person.OutstandingNotifications, message, partyOrBarID)
		}
		if person.Platform == "Android" {
			intermediateQueryResult = sendAndroidPushNotification(person.DeviceToken, message, partyOrBarID)
		}
	}
	if intermediateQueryResult.Succeeded == false {
		dynamodbCalls = append(dynamodbCalls, DynamodbCall{Succeeded: false, Error: intermediateQueryResult.Error})
	}

	// (Step 2) Creating Notification in our dynamoDB so that user can see a history of their notifications
	intermediateQueryResult = createNotification(person.FacebookID, message, partyOrBarID)
	if intermediateQueryResult.Succeeded == false {
		dynamodbCalls = append(dynamodbCalls, DynamodbCall{Succeeded: false, Error: intermediateQueryResult.Error})
	}

	queryResult.People = nil
	queryResult.DynamodbCalls = nil
	if len(dynamodbCalls) > 0 {
		queryResult.DynamodbCalls = dynamodbCalls
	}
	return queryResult
}

func createAndSendNotificationsToThesePeople(facebookIDs []string, message string, partyOrBarID string) QueryResult {
	var queryResult = QueryResult{}
	if len(facebookIDs) == 0 {
		queryResult.Succeeded = true
		return queryResult
	}
	dynamodbCalls := make([]DynamodbCall, 0)

	var intermediateQueryResult QueryResult
	// update each person's number of outstanding push notifications before we send
	//		the push notification payload so that we send the correct badger number
	for i := 0; i < len(facebookIDs); i++ {
		intermediateQueryResult = QueryResult{}
		intermediateQueryResult.Succeeded = true

		intermediateQueryResult = incrementOutstandingNotificationCountForPersonHelper(facebookIDs[i])
		if intermediateQueryResult.Succeeded == false {
			dynamodbCalls = append(dynamodbCalls, DynamodbCall{Succeeded: false, Error: intermediateQueryResult.Error})
		}
	}

	queryResult = getPeople(facebookIDs)
	if queryResult.Succeeded == false {
		return queryResult
	}

	people := queryResult.People
	for i := 0; i < len(people); i++ {
		intermediateQueryResult = QueryResult{}
		intermediateQueryResult.Succeeded = true

		person := people[i]
		// (Step 1) Sending Push Notifications
		if person.Platform != "Unknown" && person.DeviceToken != "Unknown" {
			if person.Platform == "iOS" {
				intermediateQueryResult = sendiOSPushNotification(person.DeviceToken, person.OutstandingNotifications, message, partyOrBarID)
			}
			if person.Platform == "Android" {
				intermediateQueryResult = sendAndroidPushNotification(person.DeviceToken, message, partyOrBarID)
			}
		}
		if intermediateQueryResult.Succeeded == false {
			dynamodbCalls = append(dynamodbCalls, DynamodbCall{Succeeded: false, Error: intermediateQueryResult.Error})
		}

		// (Step 2) Creating Notifications in our dynamoDB so that user's can see a history of their notifications
		intermediateQueryResult = createNotification(person.FacebookID, message, partyOrBarID)
		if intermediateQueryResult.Succeeded == false {
			dynamodbCalls = append(dynamodbCalls, DynamodbCall{Succeeded: false, Error: intermediateQueryResult.Error})
		}
	}
	queryResult.People = nil
	queryResult.DynamodbCalls = nil
	if len(dynamodbCalls) > 0 {
		queryResult.DynamodbCalls = dynamodbCalls
	}
	return queryResult
}

func testGetPeople(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	people := []string{"10155227369101712", "10216576646672295", "10154326505409816"}
	queryResult := getPeople(people)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func getPeople(facebookIDs []string) QueryResult {
	var queryResult = QueryResult{}
	if len(facebookIDs) == 0 {
		queryResult.Succeeded = true
		return queryResult
	}
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
		queryResult.Error = "getPeople function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var batchGetItemInput = dynamodb.BatchGetItemInput{}
	attributesAndValues := make([]map[string]*dynamodb.AttributeValue, len(facebookIDs))
	for i := 0; i < len(facebookIDs); i++ {
		var attributeValue = dynamodb.AttributeValue{}
		attributeValue.SetS(facebookIDs[i])
		attributesAndValues[i] = make(map[string]*dynamodb.AttributeValue)
		attributesAndValues[i]["facebookID"] = &attributeValue
	}
	var keysAndAttributes dynamodb.KeysAndAttributes
	keysAndAttributes.SetKeys(attributesAndValues)
	requestedItems := make(map[string]*dynamodb.KeysAndAttributes)
	requestedItems["Person"] = &keysAndAttributes
	batchGetItemInput.SetRequestItems(requestedItems)
	batchGetItemOutput, err2 := getter.DynamoDB.BatchGetItem(&batchGetItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getPeople function: BatchGetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := batchGetItemOutput.Responses
	people := make([]PersonData, len(facebookIDs))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data["Person"], &people)
	if jsonErr != nil {
		queryResult.Error = "getPeople function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.People = people
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func markNotificationAsSeen(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	notificationID := r.Form.Get("notificationID")
	queryResult := markNotificationAsSeenHelper(notificationID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func markNotificationAsSeenHelper(notificationID string) QueryResult {
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
		queryResult.Error = "markNotificationAsSeenHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var hasBeenSeenString = "hasBeenSeen"
	expressionAttributeNames["#hasBeenSeen"] = &hasBeenSeenString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var hasBeenSeenAttributeValue = dynamodb.AttributeValue{}
	hasBeenSeenAttributeValue.SetBOOL(true)
	expressionValuePlaceholders[":hasBeenSeen"] = &hasBeenSeenAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(notificationID)
	keyMap["notificationID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Notification")
	updateExpression := "SET #hasBeenSeen=:hasBeenSeen"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "markNotificationAsSeenHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func getNotificationsForPerson(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	facebookID := r.Form.Get("facebookID")
	queryResult := getNotificationsForPersonHelper(facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func getNotificationsForPersonHelper(facebookID string) QueryResult {
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
		queryResult.Error = "getNotificationsForPersonHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var receiverFacebookIDAttributeValue = dynamodb.AttributeValue{}
	receiverFacebookIDAttributeValue.SetS(facebookID)
	expressionValuePlaceholders[":receiverFacebookID"] = &receiverFacebookIDAttributeValue

	var queryInput = dynamodb.QueryInput{}
	queryInput.SetTableName("Notification")
	queryInput.SetIndexName("NotificationFacebookIDIndex")
	queryInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	keyConditionExpression := "receiverFacebookID = :receiverFacebookID"
	queryInput.KeyConditionExpression = &keyConditionExpression

	queryOutput, err2 := getter.DynamoDB.Query(&queryInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getNotificationsForPersonHelper function: Query error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	// Items []map[string]*AttributeValue `type:"list"`
	data := queryOutput.Items
	notifications := make([]Notification, len(data))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &notifications)

	if jsonErr != nil {
		queryResult.Error = "getNotificationsForPersonHelper function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.Notifications = notifications
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func testCreateNotification(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	receiverFacebookID := r.Form.Get("receiverFacebookID")
	message := r.Form.Get("message")
	partyOrBarID := r.Form.Get("partyOrBarID")
	queryResult := createNotification(receiverFacebookID, message, partyOrBarID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func createNotification(receiverFacebookID string, message string, partyOrBarID string) QueryResult {
	notificationID := strconv.FormatUint(getRandomID(), 10)
	hasBeenSeen := false
	timeTwoWeeksFromNow := time.Now().Add(time.Duration(336) * time.Hour) // 336 hours in two weeks
	expiresAt := strconv.FormatInt(timeTwoWeeksFromNow.Unix(), 10)

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
		queryResult.Error = "createNotification function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally

	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var notificationIDAttributeValue = dynamodb.AttributeValue{}
	var receiverFacebookIDAttributeValue = dynamodb.AttributeValue{}
	var messageAttributeValue = dynamodb.AttributeValue{}
	var partyOrBarIDAttributeValue = dynamodb.AttributeValue{}
	var hasBeenSeenAttributeValue = dynamodb.AttributeValue{}
	var expiresAtAttributeValue = dynamodb.AttributeValue{}
	notificationIDAttributeValue.SetS(notificationID)
	receiverFacebookIDAttributeValue.SetS(receiverFacebookID)
	messageAttributeValue.SetS(message)
	partyOrBarIDAttributeValue.SetS(partyOrBarID)
	hasBeenSeenAttributeValue.SetBOOL(hasBeenSeen)
	expiresAtAttributeValue.SetN(expiresAt)

	expressionValues["notificationID"] = &notificationIDAttributeValue
	expressionValues["receiverFacebookID"] = &receiverFacebookIDAttributeValue
	expressionValues["message"] = &messageAttributeValue
	expressionValues["partyOrBarID"] = &partyOrBarIDAttributeValue
	expressionValues["hasBeenSeen"] = &hasBeenSeenAttributeValue
	expressionValues["expiresAt"] = &expiresAtAttributeValue

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetTableName("Notification")
	putItemInput.SetItem(expressionValues)

	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "createNotification function: PutItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
	} else {
		queryResult.DynamodbCalls = nil
	}
	queryResult.Succeeded = true
	return queryResult
}

func testSendiOSPushNotification(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	deviceToken := r.Form.Get("deviceToken")
	queryResult := sendiOSPushNotification(deviceToken, 1, "This is a message!", "000111")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func testSendAndroidPushNotification(w http.ResponseWriter, r *http.Request) {
	queryResult := sendAndroidPushNotification("d8ueLVJ2Et0:APA91bF6RtEWZx0lnm63Q74cUTAkZbhSxreX8laylRvUh8qB4F8mSq7Gu2mmkflBOuur6JPoBKH9LnUvplVwtdtFi0fm7N3DNKVFBJV5geJsoeKZL1qHoDtmnhf0MmxopMw4j0bWsjfr455T5MiWb26ue-LBZmG5Mg", "This is a message!", "000111")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func sendiOSPushNotification(deviceToken string, outstandingNotifications uint64, message string, partyOrBarID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false

	authKey, err := token.AuthKeyFromFile("./src/AuthKey_8252N9Z82W.p8")
	if err != nil {
		queryResult.Error = err.Error()
		return queryResult
	}

	token := &token.Token{
		AuthKey: authKey,
		// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
		KeyID: "8252N9Z82W",
		// TeamID from developer account (View Account -> Membership)
		TeamID: "3SM66DY534",
	}

	outstandingNotificationCount := strconv.FormatUint(outstandingNotifications, 10)

	notification := &apns2.Notification{
		DeviceToken: deviceToken,
		Topic:       "BumpinBundleIdentifier",
		Payload:     []byte(`{"aps":{"alert":"` + message + `","badge":` + outstandingNotificationCount + `},"partyOrBarID":"` + partyOrBarID + `"}`),
	}

	client := apns2.NewTokenClient(token)
	client.Production()
	status, err2 := client.Push(notification)
	if err2 != nil {
		queryResult.Error = "Production push notification failed: " + err2.Error()
		return queryResult
	}

	if status.Reason != "" {
		client.Development()
		status2, err3 := client.Push(notification)
		if err3 != nil {
			queryResult.Error = "Development push notification failed: " + err3.Error()
			return queryResult
		}
		if status2.Reason != "" {
			queryResult.Error = "Development push notification failed: " + status2.Reason
			return queryResult
		}
	}

	queryResult.Succeeded = true
	return queryResult
}

// server key = AAAAo_YT2fc:APA91bEV1ctVnAhvWzO7uOpuMBcHpwYu1LaGDgHF3KZ4GtdY1yocH90Vc_fvFlmtGDKib1vYA24ci5QUdaoozpeI_kfd9QdHwGS2L8JNDd6AZh1I-zGZ8COLEPp75c_wlAG_iFE1NbIZ
// sender id = 704208165367
// deviceToken = f2R352-yppw:APA91bGPjrkqk0ChvcG763aROYemYkp0WvXxE5yRA0vGw0IPAu_wsfRb8wN3qSDDG_lONNRjaXl0bfYFNI1Pxr6UaX86iRxhOuwcqPZ3WfOt2N3xSgFT-_4z2AXUUcD2_CZkd3QyRqTR

func sendAndroidPushNotification(deviceToken string, message string, partyOrBarID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false

	data := map[string]string{
		"title": message,
		//"message":      "Hello!",
		"partyOrBarID": partyOrBarID,
	}

	ids := []string{
		deviceToken,
	}
	/*
	  xds := []string{
	      "token5",
	      "token6",
	      "token7",
	  }*/

	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(ids, data)
	//c.AppendDevices(xds)

	_, err := c.Send()

	if err != nil {
		queryResult.Error = err.Error()
		return queryResult
	}
	queryResult.Succeeded = true
	return queryResult
}

func testiOSPushNotification(w http.ResponseWriter, r *http.Request) {

	authKey, err := token.AuthKeyFromFile("./src/AuthKey_8252N9Z82W.p8")
	if err != nil {
		log.Fatal("token error:", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
		KeyID: "8252N9Z82W",
		// TeamID from developer account (View Account -> Membership)
		TeamID: "3SM66DY534",
	}

	notification := &apns2.Notification{
		DeviceToken: "ff63ea4106df5eb744b3289270976083f52bd0abe225d43501b214904ece3d9c",
		Topic:       "BumpinBundleIdentifier",
		Payload:     []byte(`{"aps":{"alert":"Hello!"},"partyOrBarID":"12345"}`),
	}

	client := apns2.NewTokenClient(token)
	res, err := client.Push(notification)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(res)
}

// server key = AAAAo_YT2fc:APA91bEV1ctVnAhvWzO7uOpuMBcHpwYu1LaGDgHF3KZ4GtdY1yocH90Vc_fvFlmtGDKib1vYA24ci5QUdaoozpeI_kfd9QdHwGS2L8JNDd6AZh1I-zGZ8COLEPp75c_wlAG_iFE1NbIZ
// sender id = 704208165367
// deviceToken = f2R352-yppw:APA91bGPjrkqk0ChvcG763aROYemYkp0WvXxE5yRA0vGw0IPAu_wsfRb8wN3qSDDG_lONNRjaXl0bfYFNI1Pxr6UaX86iRxhOuwcqPZ3WfOt2N3xSgFT-_4z2AXUUcD2_CZkd3QyRqTR

func testAndroidPushNotification(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"title": "Hello!",
		//"message":      "Hello!",
		"partyOrBarID": "123456",
	}

	ids := []string{
		"eEpRROIRD44:APA91bGxfaDudOSpNvggPsw-Q-DpjNtgXzh9CRFTGMKb38hDlfVscvPwcKSqMmb06Bg2vDBQgMcqpGHjE5i_l8UnlzPBVwp2_gyB0TabsU8Rc0n1KG9B4kVAk-Rbl-aoL04COhiL5pCZ",
	}
	/*
	  xds := []string{
	      "token5",
	      "token6",
	      "token7",
	  }*/

	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(ids, data)
	//c.AppendDevices(xds)

	status, err := c.Send()

	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(status)
}
