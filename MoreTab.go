package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Updates your block list
func updateBlockList(w http.ResponseWriter, r *http.Request) {
	myFacebookID := r.URL.Query().Get("myFacebookID")
	var additionsList []string
	var removalsList []string
	if r.URL.Query().Get("additionsList") != "" {
		additionsList = strings.Split(r.URL.Query().Get("additionsList"), ",")
	}
	if r.URL.Query().Get("removalsList") != "" {
		removalsList = strings.Split(r.URL.Query().Get("removalsList"), ",")
	}
	additionsQueryResult := additionsBlockListHelper(myFacebookID, additionsList)
	removalsQueryResult := removalsBlockListHelper(myFacebookID, removalsList)
	queryResult := convertTwoQueryResultsToOne(additionsQueryResult, removalsQueryResult)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func additionsBlockListHelper(myFacebookID string, additionsList []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	queryResult.DynamodbCalls = make([]DynamodbCall, len(additionsList))
	for i := 0; i < len(additionsList); i++ {
		oneOfTheQueryResults := addPersonToBlockList(myFacebookID, additionsList[i])
		queryResult.Error = queryResult.Error + oneOfTheQueryResults.Error
		queryResult.Succeeded = queryResult.Succeeded && oneOfTheQueryResults.Succeeded
		queryResult.DynamodbCalls[i] = oneOfTheQueryResults.DynamodbCalls[0]
	}
	return queryResult
}

func removalsBlockListHelper(myFacebookID string, removalsList []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	queryResult.DynamodbCalls = make([]DynamodbCall, len(removalsList))
	for i := 0; i < len(removalsList); i++ {
		oneOfTheQueryResults := removePersonFromBlockList(myFacebookID, removalsList[i])
		queryResult.Error = queryResult.Error + oneOfTheQueryResults.Error
		queryResult.Succeeded = queryResult.Succeeded && oneOfTheQueryResults.Succeeded
		queryResult.DynamodbCalls[i] = oneOfTheQueryResults.DynamodbCalls[0]
	}
	return queryResult
}

func addPersonToBlockList(myFacebookID string, theirFacebookID string) QueryResult {
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
		queryResult.Error = "addPersonToBlocklist function: session creation error. " + err.Error()
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
		dynamodbCall.Error = "addPersonToBlockList function: UpdateItem error. " + "Error adding " + theirFacebookID + " to the Block List. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}

func removePersonFromBlockList(myFacebookID string, theirFacebookID string) QueryResult {
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
		queryResult.Error = "removePersonFromBlocklist function: session creation error. " + err.Error()
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
		dynamodbCall.Error = "removePersonFromBlockList function: UpdateItem error. " + "Error removing " + theirFacebookID + " from the Block List. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}
