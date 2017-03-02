package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Find all parties I'm invited to
func myParties(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.URL.Query().Get("partyID"), ",")
	data, err := findMyParties(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	parties := make([]PartyData, len(params))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data["Party"], &parties)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(parties)
}

func findMyParties(params []string) (map[string][]map[string]*dynamodb.AttributeValue, error) {
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
	var batchGetItemInput = dynamodb.BatchGetItemInput{}
	attributesAndValues := make([]map[string]*dynamodb.AttributeValue, len(params))
	for i := 0; i < len(params); i++ {
		var attributeValue = dynamodb.AttributeValue{}
		attributeValue.SetN(params[i])
		attributesAndValues[i] = make(map[string]*dynamodb.AttributeValue)
		attributesAndValues[i]["partyID"] = &attributeValue
	}

	var keysAndAttributes dynamodb.KeysAndAttributes
	keysAndAttributes.SetKeys(attributesAndValues)

	requestedItems := make(map[string]*dynamodb.KeysAndAttributes)
	requestedItems["Party"] = &keysAndAttributes
	batchGetItemInput.SetRequestItems(requestedItems)
	//getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"partyID": &attributeValue})
	batchGetItemOutput, err2 := getter.DynamoDB.BatchGetItem(&batchGetItemInput)
	return batchGetItemOutput.Responses, err2
}

// Find all bars that I'm close to
func barsCloseToMe(w http.ResponseWriter, r *http.Request) {
	latitude, latitudeErr := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	longitude, longitudeErr := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)

	if latitudeErr != nil {
		http.Error(w, "HTTP get request latitude parameter messed up: "+latitudeErr.Error(), http.StatusInternalServerError)
		return
	}
	if longitudeErr != nil {
		http.Error(w, "HTTP get request longitude parameter messed up: "+longitudeErr.Error(), http.StatusInternalServerError)
		return
	}
	data, err := findBarsCloseToMe(latitude, longitude)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bars := make([]BarData, len(data))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &bars)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(bars)
}

func findBarsCloseToMe(latitude float64, longitude float64) ([]map[string]*dynamodb.AttributeValue, error) {
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
	var scanItemsInput = dynamodb.ScanInput{}
	scanItemsInput.SetTableName("Bar")
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	// approximately getting bars within a 100 mile radius
	// # of degrees / 100 miles = 1.45 degrees of latitude
	degreesOfLatitudeWhichEqual100Miles := 1.45
	// # of degrees / 100 miles = (1 degree of longitude / (69.1703 * COS(Latitude * 0.0174533)) ) * 100 miles
	degreesOfLongitudeWhichEqual100Miles := (1 / (69.1703 * math.Cos(latitude*0.0174533))) * 100
	latitudeSouth := latitude - degreesOfLatitudeWhichEqual100Miles
	latitudeNorth := latitude + degreesOfLatitudeWhichEqual100Miles
	longitudeEast := longitude - degreesOfLongitudeWhichEqual100Miles
	longitudeWest := longitude + degreesOfLongitudeWhichEqual100Miles
	latitudeSouthAttributeValue := dynamodb.AttributeValue{}
	latitudeNorthAttributeValue := dynamodb.AttributeValue{}
	longitudeEastAttributeValue := dynamodb.AttributeValue{}
	longitudeWestAttributeValue := dynamodb.AttributeValue{}
	latitudeSouthAttributeValue.SetN(strconv.FormatFloat(latitudeSouth, 'f', -1, 64))
	latitudeNorthAttributeValue.SetN(strconv.FormatFloat(latitudeNorth, 'f', -1, 64))
	longitudeEastAttributeValue.SetN(strconv.FormatFloat(longitudeEast, 'f', -1, 64))
	longitudeWestAttributeValue.SetN(strconv.FormatFloat(longitudeWest, 'f', -1, 64))
	expressionValuePlaceholders[":latitudeSouth"] = &latitudeSouthAttributeValue
	expressionValuePlaceholders[":latitudeNorth"] = &latitudeNorthAttributeValue
	expressionValuePlaceholders[":longitudeEast"] = &longitudeEastAttributeValue
	expressionValuePlaceholders[":longitudeWest"] = &longitudeWestAttributeValue
	scanItemsInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	scanItemsInput.SetFilterExpression("(latitude BETWEEN :latitudeSouth AND :latitudeNorth) AND (longitude BETWEEN :longitudeEast AND :longitudeWest)")
	scanItemsOutput, err2 := getter.DynamoDB.Scan(&scanItemsInput)
	return scanItemsOutput.Items, err2
}

// Change my attendance status to a party
func changeAttendanceStatusToParty(w http.ResponseWriter, r *http.Request) {
	partyID := r.URL.Query().Get("partyID")
	facebookID := r.URL.Query().Get("facebookID")
	status := r.URL.Query().Get("status")

	data, err := changeAttendanceStatusToPartyHelper(partyID, facebookID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var party = PartyData{}
	jsonErr := dynamodbattribute.UnmarshalMap(data, &party)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(party)
}

func changeAttendanceStatusToPartyHelper(partyID string, facebookID string, status string) (map[string]*dynamodb.AttributeValue, error) {
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
	invitees := "invitees"
	stringStatus := "status"
	expressionAttributeNames["#i"] = &invitees
	expressionAttributeNames["#f"] = &facebookID
	expressionAttributeNames["#s"] = &stringStatus
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var statusAttributeValue = dynamodb.AttributeValue{}
	statusAttributeValue.SetS(status)
	expressionValuePlaceholders[":status"] = &statusAttributeValue
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #i.#f.#s=:status"
	updateItemInput.UpdateExpression = &updateExpression

	updateItemOutput, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	return updateItemOutput.Attributes, err2
}

// Change my attendance status to a bar (add my info to the attendees map if need be)
func changeAttendanceStatusToBar(w http.ResponseWriter, r *http.Request) {
	barID := r.URL.Query().Get("barID")
	facebookID := r.URL.Query().Get("facebookID")
	atBar, atBarConvErr := strconv.ParseBool(r.URL.Query().Get("atBar"))
	isMale, isMaleConvErr := strconv.ParseBool(r.URL.Query().Get("isMale"))
	name := r.URL.Query().Get("name")
	rating := r.URL.Query().Get("rating")
	status := r.URL.Query().Get("status")
	timeLastRated := r.URL.Query().Get("timeLastRated")

	if atBarConvErr != nil {
		http.Error(w, "Parameter issue: "+atBarConvErr.Error(), http.StatusInternalServerError)
		return
	}
	if isMaleConvErr != nil {
		http.Error(w, "Parameter issue: "+atBarConvErr.Error(), http.StatusInternalServerError)
		return
	}

	data, err := changeAttendanceStatusToBarHelper(barID, facebookID, atBar, isMale, name, rating, status, timeLastRated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var bar = BarData{}
	jsonErr := dynamodbattribute.UnmarshalMap(data, &bar)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(bar)
}

func changeAttendanceStatusToBarHelper(barID string, facebookID string, atBar bool, isMale bool, name string, rating string, status string, timeLastRated string) (map[string]*dynamodb.AttributeValue, error) {
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
	attendees := "attendees"
	expressionAttributeNames["#a"] = &attendees
	expressionAttributeNames["#f"] = &facebookID
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)

	attendeeMap := make(map[string]*dynamodb.AttributeValue)
	var attendee = dynamodb.AttributeValue{}
	var atBarAttribute = dynamodb.AttributeValue{}
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	atBarAttribute.SetBOOL(atBar)
	isMaleAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	ratingAttribute.SetS(rating)
	statusAttribute.SetS(status)
	timeLastRatedAttribute.SetS(timeLastRated)
	attendeeMap["atBar"] = &atBarAttribute
	attendeeMap["isMale"] = &isMaleAttribute
	attendeeMap["name"] = &nameAttribute
	attendeeMap["rating"] = &ratingAttribute
	attendeeMap["status"] = &statusAttribute
	attendeeMap["timeLastRated"] = &timeLastRatedAttribute
	attendee.SetM(attendeeMap)
	expressionValuePlaceholders[":attendee"] = &attendee

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #a.#f=:attendee"
	updateItemInput.UpdateExpression = &updateExpression

	updateItemOutput, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	return updateItemOutput.Attributes, err2
}
