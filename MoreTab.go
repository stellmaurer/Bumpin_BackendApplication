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
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func decrementNumberOfFriendsThatMightGoOutForTheseFriends(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var queryResult = QueryResult{}
	// facebookID := r.Form.Get("facebookID")
	if r.Form.Get("friendFacebookIDs") == "" {
		queryResult.Succeeded = true
		queryResult.DynamodbCalls = nil
	} else {
		friendFacebookIDs := strings.Split(r.Form.Get("friendFacebookIDs"), ",")
		queryResult = decrementNumberOfFriendsThatMightGoOutForTheseFriendsHelper(friendFacebookIDs)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func decrementNumberOfFriendsThatMightGoOutForTheseFriendsHelper(facebookIDs []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if len(facebookIDs) == 0 {
		queryResult.Succeeded = true
		return queryResult
	}
	dynamodbCalls := make([]DynamodbCall, 0)

	var intermediateQueryResult QueryResult

	for i := 0; i < len(facebookIDs); i++ {
		intermediateQueryResult = QueryResult{}
		intermediateQueryResult.Succeeded = true

		intermediateQueryResult = decrementNumberOfFriendsThatMightGoOutForThisPerson(facebookIDs[i])
		if intermediateQueryResult.Succeeded == false {
			dynamodbCalls = append(dynamodbCalls, DynamodbCall{Succeeded: false, Error: intermediateQueryResult.Error})
		}
	}

	queryResult.DynamodbCalls = dynamodbCalls
	if len(dynamodbCalls) == 0 {
		queryResult.DynamodbCalls = nil
	}
	queryResult.Succeeded = true
	return queryResult
}

func decrementNumberOfFriendsThatMightGoOutForThisPerson(facebookID string) QueryResult {
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
		queryResult.Error = "decrementNumberOfFriendsThatMightGoOutForThisPerson function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var numberOfFriendsThatMightGoOutString = "numberOfFriendsThatMightGoOut"
	expressionAttributeNames["#numberOfFriendsThatMightGoOut"] = &numberOfFriendsThatMightGoOutString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var negativeOneAttributeValue = dynamodb.AttributeValue{}
	negativeOneAttributeValue.SetN("-1")
	expressionValuePlaceholders[":negativeOne"] = &negativeOneAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "ADD #numberOfFriendsThatMightGoOut :negativeOne"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "decrementNumberOfFriendsThatMightGoOutForThisPerson function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func incrementNumberOfFriendsThatMightGoOutForTheseFriends(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var queryResult = QueryResult{}
	// facebookID := r.Form.Get("facebookID")
	if r.Form.Get("friendFacebookIDs") == "" {
		queryResult.Succeeded = true
		queryResult.DynamodbCalls = nil
	} else {
		friendFacebookIDs := strings.Split(r.Form.Get("friendFacebookIDs"), ",")
		queryResult = incrementNumberOfFriendsThatMightGoOutForTheseFriendsHelper(friendFacebookIDs)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func incrementNumberOfFriendsThatMightGoOutForTheseFriendsHelper(facebookIDs []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if len(facebookIDs) == 0 {
		queryResult.Succeeded = true
		return queryResult
	}
	dynamodbCalls := make([]DynamodbCall, 0)

	var intermediateQueryResult QueryResult

	for i := 0; i < len(facebookIDs); i++ {
		intermediateQueryResult = QueryResult{}
		intermediateQueryResult.Succeeded = true

		intermediateQueryResult = incrementNumberOfFriendsThatMightGoOutForThisPerson(facebookIDs[i])
		if intermediateQueryResult.Succeeded == false {
			dynamodbCalls = append(dynamodbCalls, DynamodbCall{Succeeded: false, Error: intermediateQueryResult.Error})
		}
	}

	queryResult.DynamodbCalls = dynamodbCalls
	if len(dynamodbCalls) == 0 {
		queryResult.DynamodbCalls = nil
	}
	queryResult.Succeeded = true
	return queryResult
}

func incrementNumberOfFriendsThatMightGoOutForThisPerson(facebookID string) QueryResult {
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
		queryResult.Error = "incrementNumberOfFriendsThatMightGoOutForThisPerson function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var numberOfFriendsThatMightGoOutString = "numberOfFriendsThatMightGoOut"
	expressionAttributeNames["#numberOfFriendsThatMightGoOut"] = &numberOfFriendsThatMightGoOutString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var oneAttributeValue = dynamodb.AttributeValue{}
	oneAttributeValue.SetN("1")
	expressionValuePlaceholders[":one"] = &oneAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "ADD #numberOfFriendsThatMightGoOut :one"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "incrementNumberOfFriendsThatMightGoOutForThisPerson function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func getFriends(w http.ResponseWriter, r *http.Request) {
	var queryResult = QueryResult{}
	if r.URL.Query().Get("facebookIDs") == "" {
		queryResult.Succeeded = true
		queryResult.DynamodbCalls = nil
	} else {
		facebookIDs := strings.Split(r.URL.Query().Get("facebookIDs"), ",")
		queryResult = getFriendsHelper(facebookIDs)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func getFriendsHelper(facebookIDs []string) QueryResult {
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
		queryResult.Error = "getFriendsHelper function: session creation error. " + err.Error()
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
		dynamodbCall.Error = "getFriendsHelper function: BatchGetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := batchGetItemOutput.Responses
	friends := make([]PersonData, len(facebookIDs))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data["Person"], &friends)
	if jsonErr != nil {
		queryResult.Error = "getFriendsHelper function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.People = friends
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func updatePersonStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	facebookID := r.Form.Get("facebookID")
	goingOut := r.Form.Get("goingOut")
	timeGoingOutStatusWasSet := r.Form.Get("timeGoingOutStatusWasSet")
	manuallySet := r.Form.Get("manuallySet")
	queryResult := updatePersonStatusHelper(facebookID, goingOut, timeGoingOutStatusWasSet, manuallySet)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func updatePersonStatusHelper(facebookID string, goingOut string, timeGoingOutStatusWasSet string, manuallySet string) QueryResult {
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
		queryResult.Error = "updatePersonStatusHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var statusString = "status"
	var goingOutString = "goingOut"
	var timeGoingOutStatusWasSetString = "timeGoingOutStatusWasSet"
	var manuallySetString = "manuallySet"
	expressionAttributeNames["#status"] = &statusString
	expressionAttributeNames["#goingOut"] = &goingOutString
	expressionAttributeNames["#timeGoingOutStatusWasSet"] = &timeGoingOutStatusWasSetString
	expressionAttributeNames["#manuallySet"] = &manuallySetString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var goingOutAttributeValue = dynamodb.AttributeValue{}
	var timeGoingOutStatusWasSetAttributeValue = dynamodb.AttributeValue{}
	var manuallySetAttributeValue = dynamodb.AttributeValue{}
	goingOutAttributeValue.SetS(goingOut)
	timeGoingOutStatusWasSetAttributeValue.SetS(timeGoingOutStatusWasSet)
	manuallySetAttributeValue.SetS(manuallySet)
	expressionValuePlaceholders[":goingOut"] = &goingOutAttributeValue
	expressionValuePlaceholders[":timeGoingOutStatusWasSet"] = &timeGoingOutStatusWasSetAttributeValue
	expressionValuePlaceholders[":manuallySet"] = &manuallySetAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "SET #status.#goingOut=:goingOut, #status.#timeGoingOutStatusWasSet=:timeGoingOutStatusWasSet, #status.#manuallySet=:manuallySet"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "updatePersonStatusHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func createBug(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bugID := strconv.FormatUint(getRandomID(), 10)
	facebookID := r.Form.Get("facebookID")
	description := r.Form.Get("description")
	queryResult := createBugHelper(bugID, facebookID, description)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func createBugHelper(bugID string, facebookID string, description string) QueryResult {
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
		queryResult.Error = "createBugHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var bugIDAttributeValue = dynamodb.AttributeValue{}
	var facebookIDAttributeValue = dynamodb.AttributeValue{}
	var descriptionAttributeValue = dynamodb.AttributeValue{}
	bugIDAttributeValue.SetS(bugID)
	facebookIDAttributeValue.SetS(facebookID)
	descriptionAttributeValue.SetS(description)
	expressionValues["bugID"] = &bugIDAttributeValue
	expressionValues["facebookID"] = &facebookIDAttributeValue
	expressionValues["description"] = &descriptionAttributeValue

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetTableName("Bug")
	putItemInput.SetItem(expressionValues)
	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "createBugHelper function: PutItem error. Err msg: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func createFeatureRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	featureRequestID := strconv.FormatUint(getRandomID(), 10)
	facebookID := r.Form.Get("facebookID")
	description := r.Form.Get("description")
	queryResult := createFeatureRequestHelper(featureRequestID, facebookID, description)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func createFeatureRequestHelper(featureRequestID string, facebookID string, description string) QueryResult {
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
		queryResult.Error = "createFeatureRequestHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var featureRequestIDAttributeValue = dynamodb.AttributeValue{}
	var facebookIDAttributeValue = dynamodb.AttributeValue{}
	var descriptionAttributeValue = dynamodb.AttributeValue{}
	featureRequestIDAttributeValue.SetS(featureRequestID)
	facebookIDAttributeValue.SetS(facebookID)
	descriptionAttributeValue.SetS(description)
	expressionValues["featureRequestID"] = &featureRequestIDAttributeValue
	expressionValues["facebookID"] = &facebookIDAttributeValue
	expressionValues["description"] = &descriptionAttributeValue

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetTableName("FeatureRequest")
	putItemInput.SetItem(expressionValues)
	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "createFeatureRequestHelper function: PutItem error. Err msg: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

// Updates your activity block list. This information is actually stored
//		in the Person objects of the blocked friends, so this method updates
//		the Person objects of the friends you block your activity from, specfically
//		the "peopleBlockingTheirActivityFromMe" attribute.
func updateActivityBlockList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	myFacebookID := r.Form.Get("myFacebookID")
	var additionsList []string
	var removalsList []string
	if r.Form.Get("additionsList") != "" {
		additionsList = strings.Split(r.Form.Get("additionsList"), ",")
	}
	if r.Form.Get("removalsList") != "" {
		removalsList = strings.Split(r.Form.Get("removalsList"), ",")
	}
	additionsQueryResult := additionsActivityBlockListHelper(myFacebookID, additionsList)
	removalsQueryResult := removalsActivityBlockListHelper(myFacebookID, removalsList)
	queryResult := convertTwoQueryResultsToOne(additionsQueryResult, removalsQueryResult)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func updateIgnoreList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	myFacebookID := r.Form.Get("myFacebookID")
	var additionsList []string
	var removalsList []string
	if r.Form.Get("additionsList") != "" {
		additionsList = strings.Split(r.Form.Get("additionsList"), ",")
	}
	if r.Form.Get("removalsList") != "" {
		removalsList = strings.Split(r.Form.Get("removalsList"), ",")
	}
	queryResult := updateIgnoreListHelper(myFacebookID, additionsList, removalsList)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func additionsActivityBlockListHelper(myFacebookID string, additionsList []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	queryResult.DynamodbCalls = make([]DynamodbCall, len(additionsList))
	for i := 0; i < len(additionsList); i++ {
		oneOfTheQueryResults := addPersonToActivityBlockList(myFacebookID, additionsList[i])
		queryResult.Error = queryResult.Error + oneOfTheQueryResults.Error
		queryResult.Succeeded = queryResult.Succeeded && oneOfTheQueryResults.Succeeded
		queryResult.DynamodbCalls[i] = oneOfTheQueryResults.DynamodbCalls[0]
	}
	if queryResult.Succeeded == true {
		queryResult.DynamodbCalls = nil
	}
	return queryResult
}

func removalsActivityBlockListHelper(myFacebookID string, removalsList []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	queryResult.DynamodbCalls = make([]DynamodbCall, len(removalsList))
	for i := 0; i < len(removalsList); i++ {
		oneOfTheQueryResults := removePersonFromActivityBlockList(myFacebookID, removalsList[i])
		queryResult.Error = queryResult.Error + oneOfTheQueryResults.Error
		queryResult.Succeeded = queryResult.Succeeded && oneOfTheQueryResults.Succeeded
		queryResult.DynamodbCalls[i] = oneOfTheQueryResults.DynamodbCalls[0]
	}
	if queryResult.Succeeded == true {
		queryResult.DynamodbCalls = nil
	}
	return queryResult
}

func addPersonToActivityBlockList(myFacebookID string, theirFacebookID string) QueryResult {
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
		queryResult.Error = "addPersonToActivityBlockList function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	peopleBlockingTheirActivityFromMe := "peopleBlockingTheirActivityFromMe"
	expressionAttributeNames["#p"] = &peopleBlockingTheirActivityFromMe
	expressionAttributeNames["#f"] = &myFacebookID

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var myFacebookIDBoolAttribute = dynamodb.AttributeValue{}
	myFacebookIDBoolAttribute.SetBOOL(true)
	expressionValuePlaceholders[":bool"] = &myFacebookIDBoolAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(theirFacebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "SET #p.#f=:bool"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "addPersonToActivityBlockList function: UpdateItem error. " + "Error adding " + theirFacebookID + " to the Block List. This person most likely does not exist anymore." + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}

func removePersonFromActivityBlockList(myFacebookID string, theirFacebookID string) QueryResult {
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
		queryResult.Error = "removePersonFromActivityBlockList function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	peopleBlockingTheirActivityFromMe := "peopleBlockingTheirActivityFromMe"
	expressionAttributeNames["#p"] = &peopleBlockingTheirActivityFromMe
	expressionAttributeNames["#f"] = &myFacebookID

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(theirFacebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "Remove #p.#f"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "removePersonFromActivityBlockList function: UpdateItem error. " + "Error removing " + theirFacebookID + " from the Block List. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}

func updateIgnoreListHelper(myFacebookID string, additionsList []string, removalsList []string) QueryResult {
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
		queryResult.Error = "updateIgnoreList function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	updateExpression := ""

	peopleToIgnoreString := "peopleToIgnore"
	expressionAttributeNames["#peopleToIgnore"] = &peopleToIgnoreString
	// Update Expression will look like this:
	//   "SET #peopleToIgnore.#facebookID1=:bool, #peopleToIgnore.#facebookID2=:bool Remove #peopleToIgnore.#facebookID3, #peopleToIgnore.#facebookID4"
	if len(additionsList) >= 1 {
		fmt.Println(additionsList)
		// Update attribute names
		for i := 0; i < len(additionsList); i++ {
			expressionAttributeNames["#"+additionsList[i]] = &additionsList[i]
		}
		// Update value placeholders
		var boolAttributeValue = dynamodb.AttributeValue{}
		boolAttributeValue.SetBOOL(true)
		expressionValuePlaceholders[":bool"] = &boolAttributeValue
		// Update the update expression
		updateExpression += "SET #peopleToIgnore.#" + additionsList[0] + "=:bool"
		for i := 1; i < len(additionsList); i++ {
			updateExpression += ", " + "#peopleToIgnore.#" + additionsList[i] + "=:bool"
		}
	}
	if len(removalsList) >= 1 {
		fmt.Println(removalsList)
		// Update attribute names
		for i := 0; i < len(removalsList); i++ {
			expressionAttributeNames["#"+removalsList[i]] = &removalsList[i]
		}
		// update the update expression
		updateExpression += " REMOVE #peopleToIgnore.#" + removalsList[0]
		for i := 1; i < len(removalsList); i++ {
			updateExpression += ", " + "#peopleToIgnore.#" + removalsList[i]
		}
	}

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(myFacebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	if len(additionsList) >= 1 {
		updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	}
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "updateIgnoreList function: UpdateItem error. A person in the additions list most likely does not exist anymore." + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}
