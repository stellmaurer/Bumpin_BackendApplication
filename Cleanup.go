package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Deletes parties that have expired. A party expires 2 hours after it's endTime.
func deletePartiesThatHaveExpired(w http.ResponseWriter, r *http.Request) {
	queryResult, expiredParties := findPartiesThatHaveExpired()
	if queryResult.Succeeded == true {
		for i := 0; i < len(expiredParties); i++ {
			var deletePartyQueryResult = deletePartyHelper(expiredParties[i].PartyID)
			if deletePartyQueryResult.Succeeded == false {
				queryResult = convertTwoQueryResultsToOne(queryResult, deletePartyQueryResult)
			}
		}
	}
	if queryResult.Succeeded == true {
		queryResult.DynamodbCalls = nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

/*
func deletePartiesThatHaveExpiredHelper(partyIDs []string) QueryResult {
	// 2017-03-16T09:00:00-Z
}*/

func findPartiesThatHaveExpired() (QueryResult, []PartyID) {
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
		queryResult.Error = "findPartiesThatHaveExpired function: session creation error. " + err.Error()
		return queryResult, nil
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	// expired if (partyEndTime + 2 hours <= current time), which is the same as:
	//    expired if (partyEndTime <= current time - 2 hours)
	twoHoursAgo := time.Now().Add(-time.Duration(2) * time.Hour).UTC().Format("2006-01-02T15:04:05Z")

	var scanItemsInput = dynamodb.ScanInput{}

	expressionAttributeNames := make(map[string]*string)
	var endTime = "endTime"
	expressionAttributeNames["#endTime"] = &endTime

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	twoHoursAgoAttributeValue := dynamodb.AttributeValue{}
	twoHoursAgoAttributeValue.SetS(twoHoursAgo)
	expressionValuePlaceholders[":twoHoursAgo"] = &twoHoursAgoAttributeValue
	scanItemsInput.SetTableName("Party")
	scanItemsInput.SetExpressionAttributeNames(expressionAttributeNames)
	scanItemsInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	scanItemsInput.SetFilterExpression("#endTime <= :twoHoursAgo")
	scanItemsOutput, err2 := getter.DynamoDB.Scan(&scanItemsInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "findPartiesThatHaveExpired function: Scan error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult, nil
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := scanItemsOutput.Items
	parties := make([]PartyID, len(data))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &parties)
	if jsonErr != nil {
		queryResult.Error = "findPartiesThatHaveExpired function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult, nil
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult, parties
}

// PartyID : a partyID
type PartyID struct {
	PartyID string `json:"partyID"`
}
