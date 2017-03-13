package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

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
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}
