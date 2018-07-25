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

// Create a bar key for a bar owner
func createClaimKeyForBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	queryResult := createClaimKeyForBarHelper(barID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func createClaimKeyForBarHelper(barID string) QueryResult {
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
		queryResult.Error = "createClaimKeyForBar function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var claimKeyAttributeValue = dynamodb.AttributeValue{}
	var barIDAttributeValue = dynamodb.AttributeValue{}
	var key = getRandomBarKey()
	claimKeyAttributeValue.SetS(key)
	barIDAttributeValue.SetS(barID)
	expressionValues["key"] = &claimKeyAttributeValue
	expressionValues["barID"] = &barIDAttributeValue

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetTableName("BarKey")
	putItemInput.SetItem(expressionValues)
	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "createClaimKeyForBar function: PutItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Error = key
	queryResult.Succeeded = true
	return queryResult
}

// Create a bar key for a bar owner
func createBarKeyForAddress(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	address := r.Form.Get("address")
	queryResult := createBarKeyForAddressHelper(address)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func createBarKeyForAddressHelper(address string) QueryResult {
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
		queryResult.Error = "createBarKeyForAddress function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var barKeyAttributeValue = dynamodb.AttributeValue{}
	var addressAttributeValue = dynamodb.AttributeValue{}
	var key = getRandomBarKey()
	barKeyAttributeValue.SetS(key)
	addressAttributeValue.SetS(address)
	expressionValues["key"] = &barKeyAttributeValue
	expressionValues["address"] = &addressAttributeValue

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetTableName("BarKey")
	putItemInput.SetItem(expressionValues)
	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "createBarKeyForAddress function: PutItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Error = key
	queryResult.Succeeded = true
	return queryResult
}
