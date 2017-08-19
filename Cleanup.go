package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

/*
// This map maps each time zone to the Zulu time hour that it's bar attendee
//			list is scheduled to be cleaned up for.
// Math: 6 = 6 AM local time (this makes it so bar attendee lists get cleaned up from 6-7 AM)
//			 Subtract hours if it's UTC+, Add hours if it's UTC-
//			 Mod by 24 to get to Zulu (UTC) time.
//			 If the result has a decimal, round up so that cleanup is during 7-8 AM.
var zuluHourOfBarClose = map[int]float32{
	0:  16,    // 6 - 14 % 24  			  // LINT/TOST       Pacific UTC +14
	1:  17,    // 6 - 13:45 % 24 		  // CHADT           Pacific UTC +13:45
	2:  17,    // 6 - 13 % 24				  // NZDT/WST/FJST   Pacific, Antarctica UTC +13
	3:  18,    // 6 - 12:45 % 24			// CHAST           Pacific UTC +12:45
	4:  18,    // 6 - 12 % 24				  // PETT/NZST/FJT   Asia, Pacific, Antarctica UTC +12
	5:  19,    // 6 - 11 % 24				  // AEDT/VUT/MAGT   Asia, Australia, Pacific UTC +11
	6:  20,    // 6 - 10:30 % 24			// LHST/ACDT       Australia UTC +10:30
	7:  20,    // 6 - 10 % 24				  // AEST/PGT/YAKST  Australia, Antarctica, Asia, Pacific UTC +10
	8:  21,    // 6 - 9:30 % 24			  // ACST            Australia UTC +9:30
	9:  21,    // 6 - 9 % 24					// KST/JST/AWDT    Australia, Asia UTC +9
	10: 22,    // 6 - 8:45 % 24			  // ACWST           Australia UTC +8:45
	11: 22,    // 6 - 8:30 % 24			  // PYT             Asia UTC +8:30
	12: 22,    // 6 - 8 % 24					// HKT/CST/AWST    Asia, Australia, Antarctica UTC +8
	13: 23,    // 6 - 7 % 24					// WIB/CXT/HOVT    Asia, Australia, Antarctica UTC +7
	14: 0,     // 6 - 6:30 % 24			  // CCT/MMT         Asia, Indian Ocean UTC +6:30
	15: 0,     // 6 - 6 % 24					// BST/IOT/VOST    Asia, Indian Ocean, Antarctica UTC +6
	16: 1,     // 6 - 5:45 % 24			  // NPT             Asia UTC +5:45
	17: 1,     // 6 - 5:30 % 24			  // IST             Asia UTC +5:30
	18: 1,     // 6 - 5 % 24				 	// PKT/TFT/UZT     Asia, Indian Ocean, Antarctica UTC +5
	19: 2,     // 6 - 4:30 % 24			  // IRDT/AFT        Asia UTC +4:30
	20: 2,     // 6 - 4 % 24					// MUT/MSD/GET     Asia, Africa, Europe UTC +4
	21: 3,     // 6 - 3:30 % 24			  // IRST            Asia UTC +3:30
	22: 3,     // 6 - 3 % 24					// MCK/CAT/EAT     Europe, Africa, Asia, Indian Ocean, Antarctica UTC +3
	23: 4,     // 6 - 2 % 24					// CEST/EET/SAST   Europe, Africa, Asia, Antarctica UTC +2
	24: 5,     // 6 - 1 % 24					// WAT/WEST/CET    Africa, Europe, Antarctica UTC +1
	25: 6,     // 6 - 0 % 24					// UTC/GMT/Z       Africa, Antarctica, North America UTC +0
	26: 7,     // 6 + 1 % 24					// EGT/AZOT/CVT    North America, Africa UTC -1
	27: 8,     // 6 + 2 % 24					// BST/WGST        North America, South America UTC -2
	28: 9,     // 6 + 2:30 % 24			  // NDT             North America UTC -2:30
	29: 9,     // 6 + 3 % 24					// BT/WGT          North America, South America UTC -3
	30: 10,    // 6 + 3:30 % 24			  // NST             North America UTC -3:30
	31: 10,    // 6 + 4 % 24					// EDT/PYT/VET     North America, South America UTC -4
	32: 11,    // 6 + 5 % 24					// PET/CDT/EST     North America, South America UTC -5
	33: 12,    // 6 + 6 % 24					// CST/MDT         North America UTC -6
	34: 13,    // 6 + 7 % 24					// PDT/MST         North America UTC -7
	35: 14,    // 6 + 8 % 24					// PST             North America, Pacific UTC -8
	36: 15,    // 6 + 9 % 24					// HADT            North America UTC -9
	37: 16,    // 6 + 9:30 % 24			  // MART            Pacific UTC -9:30
	38: 16,    // 6 + 10 % 24				  // HAST/TAHT       Pacific UTC -10
	39: 17,    // 6 + 11 % 24				  // SST             Pacific UTC -11
	40: 18,    // 6 + 12 % 24				  // AoE             Pacific UTC -12
}
*/

// Go through all the bars in the database and find the bars that just closed.
//		Remove all attendees from a bar's attendee list if it just closed.
// Implementation idea:
//		In order to make this ^ work, we only need the bar's time zone.
//		For the sake of simplicity, we'll say that all bars close at 6 AM
//		in their time zone. We'll have attendee cleanups run every hour. This means
//		that at some point from 6 AM - 7 AM in the bar's time zone, the attendee
//		list will be cleaned up.
func cleanUpAttendeesMapForBarsThatRecentlyClosed(w http.ResponseWriter, r *http.Request) {
	// TODO: WE NEED 2 additional attributes on Bar -> timeZone (in int map key format), and zuluHourOfBarClose
	queryResult := cleanUpAttendeesMapForBarsThatRecentlyClosedHelper()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

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

// BarID : a barID
type BarID struct {
	BarID string `json:"barID"`
}

// Starts the bars that recently closed with a new fresh attendee list.
func cleanUpAttendeesMapForBarsThatRecentlyClosedHelper() QueryResult {
	queryResult, recentlyClosedBars := findBarsThatRecentlyClosed()
	fmt.Println(recentlyClosedBars)
	if queryResult.Succeeded == true {
		for i := 0; i < len(recentlyClosedBars); i++ {
			var cleanUpAttendeesMapOfBarQueryResult = cleanUpAttendeesMapOfBar(recentlyClosedBars[i].BarID)
			if cleanUpAttendeesMapOfBarQueryResult.Succeeded == false {
				queryResult = convertTwoQueryResultsToOne(queryResult, cleanUpAttendeesMapOfBarQueryResult)
			}
		}
	}
	if queryResult.Succeeded == true {
		queryResult.DynamodbCalls = nil
	}
	return queryResult
}

func findBarsThatRecentlyClosed() (QueryResult, []BarID) {
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
		queryResult.Error = "findBarsThatRecentlyClosed function: session creation error. " + err.Error()
		return queryResult, nil
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	currentTimeInZulu := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	zuluHourOfCurrentTimeString := currentTimeInZulu[11:13]

	var scanItemsInput = dynamodb.ScanInput{}

	expressionAttributeNames := make(map[string]*string)
	var attendeesMapCleanUpHourInZuluString = "attendeesMapCleanUpHourInZulu"
	expressionAttributeNames["#attendeesMapCleanUpHourInZulu"] = &attendeesMapCleanUpHourInZuluString

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	zuluHourOfCurrentTimeAttributeValue := dynamodb.AttributeValue{}
	zuluHourOfCurrentTimeAttributeValue.SetN(zuluHourOfCurrentTimeString)
	expressionValuePlaceholders[":zuluHourOfCurrentTime"] = &zuluHourOfCurrentTimeAttributeValue
	scanItemsInput.SetTableName("Bar")
	scanItemsInput.SetExpressionAttributeNames(expressionAttributeNames)
	scanItemsInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	scanItemsInput.SetFilterExpression("#attendeesMapCleanUpHourInZulu = :zuluHourOfCurrentTime")
	scanItemsOutput, err2 := getter.DynamoDB.Scan(&scanItemsInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "findBarsThatRecentlyClosed function: Scan error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult, nil
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := scanItemsOutput.Items
	bars := make([]BarID, len(data))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &bars)
	if jsonErr != nil {
		queryResult.Error = "findBarsThatRecentlyClosed function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult, nil
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult, bars
}

func cleanUpAttendeesMapOfBar(barID string) QueryResult {
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
		queryResult.Error = "cleanUpAttendeesMapOfBar function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	attendeesString := "attendees"
	expressionAttributeNames["#a"] = &attendeesString

	// Make the attendeesMap an empty map (this cleans up all the attendees because it's a new day)
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	attendeesMap := make(map[string]*dynamodb.AttributeValue)
	var attendees = dynamodb.AttributeValue{}
	attendees.SetM(attendeesMap)
	expressionValuePlaceholders[":attendees"] = &attendees

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #a=:attendees"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "cleanUpAttendeesMapOfBar function: UpdateItem error. (probable cause: this bar may not exist)" + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}
