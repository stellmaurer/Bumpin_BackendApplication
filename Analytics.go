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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func logError(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ID := r.Form.Get("ID")
	errorType := r.Form.Get("errorType")
	errorDescription := r.Form.Get("errorDescription")
	queryResult := logErrorHelper(ID, errorType, errorDescription)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func logErrorHelper(ID string, errorType string, errorDescription string) QueryResult {
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
		queryResult.Error = "logErrorHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	stringErrorType := "errorType"
	stringErrors := "errors"
	stringNumberOfErrors := "numberOfErrors"
	expressionAttributeNames["#errorTypeMap"] = &stringErrorType
	expressionAttributeNames["#errorType"] = &errorType
	expressionAttributeNames["#errorsList"] = &stringErrors
	expressionAttributeNames["#numberOfErrors"] = &stringNumberOfErrors

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var errorDescriptionAttributeValue = dynamodb.AttributeValue{}
	errorDescriptionAttributeValue.SetS(errorDescription)
	errorsList := []*dynamodb.AttributeValue{&errorDescriptionAttributeValue}
	var errorsListAttributeValue = dynamodb.AttributeValue{}
	errorsListAttributeValue.SetL(errorsList)
	expressionValuePlaceholders[":errorDescription"] = &errorsListAttributeValue
	var incrementAttributeValue = dynamodb.AttributeValue{}
	incrementAttributeValue.SetN("1")
	expressionValuePlaceholders[":increment"] = &incrementAttributeValue

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(ID)
	keyMap["ID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Analytics")

	updateExpression := "SET #errorTypeMap.#errorType.#errorsList=list_append(:errorDescription, #errorTypeMap.#errorType.#errorsList) ADD #numberOfErrors :increment"
	updateItemInput.UpdateExpression = &updateExpression
	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "logErrorHelper function: UpdateItem error. Err msg: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}
