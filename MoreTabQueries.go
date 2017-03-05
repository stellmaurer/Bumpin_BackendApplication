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
	var queryResult = QueryResult{}
	queryResult.Succeeded = additionsQueryResult.Succeeded && removalsQueryResult.Succeeded
	for i := 0; i < len(additionsQueryResult.Errors); i++ {
		queryResult.Errors = append(queryResult.Errors, additionsQueryResult.Errors[i])
	}
	for i := len(additionsQueryResult.Errors); i < len(additionsQueryResult.Errors)+len(removalsQueryResult.Errors); i++ {
		queryResult.Errors = append(queryResult.Errors, removalsQueryResult.Errors[i])
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func additionsBlockListHelper(myFacebookID string, additionsList []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	for i := 0; i < len(additionsList); i++ {
		err := addPersonToBlockList(myFacebookID, additionsList[i])
		if err != nil {
			queryResult.Errors = append(queryResult.Errors, "Error adding "+additionsList[i]+" to the Block List. "+err.Error())
			queryResult.Succeeded = false
		}
	}
	return queryResult
}

func removalsBlockListHelper(myFacebookID string, removalsList []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	for i := 0; i < len(removalsList); i++ {
		err := removePersonFromBlockList(myFacebookID, removalsList[i])
		if err != nil {
			queryResult.Errors = append(queryResult.Errors, "Error removing "+removalsList[i]+" from the Block List. "+err.Error())
			queryResult.Succeeded = false
		}
	}
	return queryResult
}

func addPersonToBlockList(myFacebookID string, theirFacebookID string) error {
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
	return err2
}

func removePersonFromBlockList(myFacebookID string, theirFacebookID string) error {
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
	return err2
}
