package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Create or update a Person in the database
func createOrUpdatePerson(w http.ResponseWriter, r *http.Request) {
	facebookID := r.URL.Query().Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.URL.Query().Get("isMale"))
	name := r.URL.Query().Get("name")

	if isMaleConvErr != nil {
		http.Error(w, "Parameter issue: "+isMaleConvErr.Error(), http.StatusInternalServerError)
		return
	}

	queryResult1 := createPersonHelper(facebookID, isMale, name)
	queryResult2 := updatePersonHelper(facebookID, isMale, name)
	var queryResult = QueryResult{}
	queryResult.Succeeded = queryResult1.Succeeded && queryResult2.Succeeded
	for i := 0; i < len(queryResult1.Errors); i++ {
		queryResult.Errors = append(queryResult.Errors, queryResult1.Errors[i])
	}
	for i := len(queryResult1.Errors); i < len(queryResult1.Errors)+len(queryResult2.Errors); i++ {
		queryResult.Errors = append(queryResult.Errors, queryResult2.Errors[i])
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult1)
}

func createPersonHelper(facebookID string, isMale bool, name string) QueryResult {
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
	var facebookIDString = "facebookID"
	expressionAttributeNames["#fbid"] = &facebookIDString

	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var barHostForAttributeValue = dynamodb.AttributeValue{}
	var facebookIDAttributeValue = dynamodb.AttributeValue{}
	var invitedToAttributeValue = dynamodb.AttributeValue{}
	var isMaleAttributeValue = dynamodb.AttributeValue{}
	var nameAttributeValue = dynamodb.AttributeValue{}
	var partyHostForAttributeValue = dynamodb.AttributeValue{}
	var peopleBlockingTheirActivityFromMeAttributeValue = dynamodb.AttributeValue{}
	barHostForAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	facebookIDAttributeValue.SetS(facebookID)
	invitedToAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	isMaleAttributeValue.SetBOOL(isMale)
	nameAttributeValue.SetS(name)
	partyHostForAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	peopleBlockingTheirActivityFromMeAttributeValue.SetM(make(map[string]*dynamodb.AttributeValue))
	expressionValues["barHostFor"] = &barHostForAttributeValue
	expressionValues["facebookID"] = &facebookIDAttributeValue
	expressionValues["invitedTo"] = &invitedToAttributeValue
	expressionValues["isMale"] = &isMaleAttributeValue
	expressionValues["name"] = &nameAttributeValue
	expressionValues["partyHostFor"] = &partyHostForAttributeValue
	expressionValues["peopleBlockingTheirActivityFromMe"] = &peopleBlockingTheirActivityFromMeAttributeValue

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetConditionExpression("attribute_not_exists(#fbid)")
	putItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	putItemInput.SetTableName("Person")
	putItemInput.SetItem(expressionValues)

	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	if err2 != nil {
		queryResult.Succeeded = false
		queryResult.Errors = append(queryResult.Errors, err2.Error())
	}
	return queryResult
}

func updatePersonHelper(facebookID string, isMale bool, name string) QueryResult {
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
	var isMaleString = "isMale"
	var nameString = "name"
	expressionAttributeNames["#isMale"] = &isMaleString
	expressionAttributeNames["#name"] = &nameString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var isMaleAttributeValue = dynamodb.AttributeValue{}
	var nameAttributeValue = dynamodb.AttributeValue{}
	isMaleAttributeValue.SetBOOL(isMale)
	nameAttributeValue.SetS(name)
	expressionValuePlaceholders[":isMale"] = &isMaleAttributeValue
	expressionValuePlaceholders[":name"] = &nameAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "SET #isMale=:isMale, #name=:name"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var queryResult = QueryResult{}
	queryResult.Succeeded = true
	if err2 != nil {
		queryResult.Succeeded = false
		queryResult.Errors = append(queryResult.Errors, err2.Error())
	}
	return queryResult
}
