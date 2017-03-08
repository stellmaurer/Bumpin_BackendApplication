package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const null = "null"

// Create a party
func createParty(w http.ResponseWriter, r *http.Request) {
	facebookID := r.URL.Query().Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.URL.Query().Get("isMale"))
	name := r.URL.Query().Get("name")
	addressLine1 := r.URL.Query().Get("addressLine1")
	addressLine2 := r.URL.Query().Get("addressLine2")
	city := r.URL.Query().Get("city")
	country := r.URL.Query().Get("country")
	details := r.URL.Query().Get("details")
	drinksProvided, drinksProvidedConvErr := strconv.ParseBool(r.URL.Query().Get("drinksProvided"))
	endTime := r.URL.Query().Get("endTime")
	feeForDrinks, feeForDrinksConvErr := strconv.ParseBool(r.URL.Query().Get("feeForDrinks"))
	invitesForNewInvitees := r.URL.Query().Get("invitesForNewInvitees")
	latitude := r.URL.Query().Get("latitude")
	longitude := r.URL.Query().Get("longitude")
	partyID := strconv.FormatUint(getRandomID(), 10)
	startTime := r.URL.Query().Get("startTime")
	stateProvinceRegion := r.URL.Query().Get("stateProvinceRegion")
	title := r.URL.Query().Get("title")
	zipCode := r.URL.Query().Get("zipCode")
	var queryResult = QueryResult{}
	if isMaleConvErr != nil {
		queryResult.Error = "createParty function: isMale parameter issue. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if drinksProvidedConvErr != nil {
		queryResult.Error = "createParty function: drinksProvided parameter issue. " + drinksProvidedConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if feeForDrinksConvErr != nil {
		queryResult.Error = "createParty function: feeForDrinks parameter issue. " + drinksProvidedConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if addressLine2 == "" {
		addressLine2 = null
	}
	if details == "" {
		details = null
	}
	queryResult = createPartyHelper(facebookID, isMale, name, addressLine1, addressLine2, city, country, details, drinksProvided, endTime, feeForDrinks, invitesForNewInvitees, latitude, longitude, partyID, startTime, stateProvinceRegion, title, zipCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

// Delete a Person from the database
func deleteParty(w http.ResponseWriter, r *http.Request) {
	partyID := r.URL.Query().Get("partyID")
	queryResult := deletePartyHelper(partyID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

// Update a party's information
func updateParty(w http.ResponseWriter, r *http.Request) {
	addressLine1 := r.URL.Query().Get("addressLine1")
	addressLine2 := r.URL.Query().Get("addressLine2")
	city := r.URL.Query().Get("city")
	country := r.URL.Query().Get("country")
	details := r.URL.Query().Get("details")
	drinksProvided, drinksProvidedConvErr := strconv.ParseBool(r.URL.Query().Get("drinksProvided"))
	endTime := r.URL.Query().Get("endTime")
	feeForDrinks, feeForDrinksConvErr := strconv.ParseBool(r.URL.Query().Get("feeForDrinks"))
	invitesForNewInvitees := r.URL.Query().Get("invitesForNewInvitees")
	latitude := r.URL.Query().Get("latitude")
	longitude := r.URL.Query().Get("longitude")
	partyID := r.URL.Query().Get("partyID")
	startTime := r.URL.Query().Get("startTime")
	stateProvinceRegion := r.URL.Query().Get("stateProvinceRegion")
	title := r.URL.Query().Get("title")
	zipCode := r.URL.Query().Get("zipCode")
	var queryResult = QueryResult{}
	if drinksProvidedConvErr != nil {
		queryResult.Error = "updateParty function: drinksProvided parameter issue. " + drinksProvidedConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if feeForDrinksConvErr != nil {
		queryResult.Error = "updateParty function: feeForDrinks parameter issue. " + drinksProvidedConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if addressLine2 == "" {
		addressLine2 = null
	}
	if details == "" {
		details = null
	}
	queryResult = updatePartyHelper(addressLine1, addressLine2, city, country, details, drinksProvided, endTime, feeForDrinks, invitesForNewInvitees, latitude, longitude, partyID, startTime, stateProvinceRegion, title, zipCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func createPartyHelper(facebookID string, isMale bool, name string, addressLine1 string, addressLine2 string, city string, country string, details string, drinksProvided bool, endTime string, feeForDrinks bool, invitesForNewInvitees string, latitude string, longitude string, partyID string, startTime string, stateProvinceRegion string, title string, zipCode string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 2)
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "updatePartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var addressLine1AttributeValue = dynamodb.AttributeValue{}
	var addressLine2AttributeValue = dynamodb.AttributeValue{}
	var cityAttributeValue = dynamodb.AttributeValue{}
	var countryAttributeValue = dynamodb.AttributeValue{}
	var detailsAttributeValue = dynamodb.AttributeValue{}
	var drinksProvidedAttributeValue = dynamodb.AttributeValue{}
	var endTimeAttributeValue = dynamodb.AttributeValue{}
	var feeForDrinksAttributeValue = dynamodb.AttributeValue{}
	var invitesForNewInviteesAttributeValue = dynamodb.AttributeValue{}
	var latitudeAttributeValue = dynamodb.AttributeValue{}
	var longitudeAttributeValue = dynamodb.AttributeValue{}
	var partyIDAttributeValue = dynamodb.AttributeValue{}
	var startTimeAttributeValue = dynamodb.AttributeValue{}
	var stateProvinceRegionAttributeValue = dynamodb.AttributeValue{}
	var titleAttributeValue = dynamodb.AttributeValue{}
	var zipCodeAttributeValue = dynamodb.AttributeValue{}
	invitesForNewInviteesAttributeValue.SetN(invitesForNewInvitees)
	addressLine1AttributeValue.SetS(addressLine1)
	addressLine2AttributeValue.SetS(addressLine2)
	cityAttributeValue.SetS(city)
	countryAttributeValue.SetS(country)
	detailsAttributeValue.SetS(details)
	drinksProvidedAttributeValue.SetBOOL(drinksProvided)
	endTimeAttributeValue.SetS(endTime)
	feeForDrinksAttributeValue.SetBOOL(feeForDrinks)
	invitesForNewInviteesAttributeValue.SetN(invitesForNewInvitees)
	latitudeAttributeValue.SetN(latitude)
	longitudeAttributeValue.SetN(longitude)
	partyIDAttributeValue.SetN(partyID)
	startTimeAttributeValue.SetS(startTime)
	stateProvinceRegionAttributeValue.SetS(stateProvinceRegion)
	titleAttributeValue.SetS(title)
	zipCodeAttributeValue.SetN(zipCode)
	expressionValues["addressLine1"] = &addressLine1AttributeValue
	expressionValues["addressLine2"] = &addressLine2AttributeValue
	expressionValues["city"] = &cityAttributeValue
	expressionValues["country"] = &countryAttributeValue
	expressionValues["details"] = &detailsAttributeValue
	expressionValues["drinksProvided"] = &drinksProvidedAttributeValue
	expressionValues["endTime"] = &endTimeAttributeValue
	expressionValues["feeForDrinks"] = &feeForDrinksAttributeValue
	expressionValues["invitesForNewInvitees"] = &invitesForNewInviteesAttributeValue
	expressionValues["latitude"] = &latitudeAttributeValue
	expressionValues["longitude"] = &longitudeAttributeValue
	expressionValues["partyID"] = &partyIDAttributeValue
	expressionValues["startTime"] = &startTimeAttributeValue
	expressionValues["stateProvinceRegion"] = &stateProvinceRegionAttributeValue
	expressionValues["title"] = &titleAttributeValue
	expressionValues["zipCode"] = &zipCodeAttributeValue

	// set yourself as an invitee to your own party so that you can rate it
	inviteesMap := make(map[string]*dynamodb.AttributeValue)
	var invitees = dynamodb.AttributeValue{}
	inviteeMap := make(map[string]*dynamodb.AttributeValue)
	var invitee = dynamodb.AttributeValue{}
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var numberOfInvitationsLeftAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	isMaleAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	numberOfInvitationsLeftAttribute.SetN("0")
	ratingAttribute.SetS("N")
	statusAttribute.SetS("G")
	timeLastRatedAttribute.SetS("01/01/2000 00:00:00")
	inviteeMap["isMale"] = &isMaleAttribute
	inviteeMap["name"] = &nameAttribute
	inviteeMap["numberOfInvitationsLeft"] = &numberOfInvitationsLeftAttribute
	inviteeMap["rating"] = &ratingAttribute
	inviteeMap["status"] = &statusAttribute
	inviteeMap["timeLastRated"] = &timeLastRatedAttribute
	invitee.SetM(inviteeMap)
	inviteesMap[facebookID] = &invitee
	invitees.SetM(inviteesMap)
	expressionValues["invitees"] = &invitees

	hostsMap := make(map[string]*dynamodb.AttributeValue)
	var hosts = dynamodb.AttributeValue{}
	hostMap := make(map[string]*dynamodb.AttributeValue)
	var host = dynamodb.AttributeValue{}
	var isMainHostAttribute = dynamodb.AttributeValue{}
	isMainHostAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	hostMap["isMainHost"] = &isMainHostAttribute
	hostMap["name"] = &nameAttribute
	host.SetM(hostMap)
	hostsMap[facebookID] = &host
	hosts.SetM(hostsMap)
	expressionValues["hosts"] = &hosts

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetTableName("Party")
	putItemInput.SetItem(expressionValues)
	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "createPartyHelper function: PutItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall

	// Now we need to update the person's information to let them
	//     know that they are invited to this party and that they
	//     are hosting this party.
	expressionAttributeNames := make(map[string]*string)
	invitedTo := "invitedTo"
	partyHostFor := "partyHostFor"
	expressionAttributeNames["#invitedTo"] = &invitedTo
	expressionAttributeNames["#partyHostFor"] = &partyHostFor
	expressionAttributeNames["#partyID"] = &partyID
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var partyIDBoolAttribute = dynamodb.AttributeValue{}
	partyIDBoolAttribute.SetBOOL(true)
	expressionValuePlaceholders[":bool"] = &partyIDBoolAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "SET #invitedTo.#partyID=:bool, #partyHostFor.#partyID=:bool"
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall2 = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall2.Error = "createPartyHelper function: UpdateItem error (probable cause: your facebookID isn't in the database). " + updateItemOutputErr.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		return queryResult
	}
	dynamodbCall2.Succeeded = true
	queryResult.DynamodbCalls[1] = dynamodbCall2
	queryResult.Succeeded = true
	return queryResult
}

func updatePartyHelper(addressLine1 string, addressLine2 string, city string, country string, details string, drinksProvided bool, endTime string, feeForDrinks bool, invitesForNewInvitees string, latitude string, longitude string, partyID string, startTime string, stateProvinceRegion string, title string, zipCode string) QueryResult {
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
		queryResult.Error = "updatePartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var addressLine1String = "addressLine1"
	var addressLine2String = "addressLine2"
	var cityString = "city"
	var countryString = "country"
	var detailsString = "details"
	var drinksProvidedString = "drinksProvided"
	var endTimeString = "endTime"
	var feeForDrinksString = "feeForDrinks"
	var invitesForNewInviteesString = "invitesForNewInvitees"
	var latitudeString = "latitude"
	var longitudeString = "longitude"
	var startTimeString = "startTime"
	var stateProvinceRegionString = "stateProvinceRegion"
	var titleString = "title"
	var zipCodeString = "zipCode"
	expressionAttributeNames["#addressLine1"] = &addressLine1String
	expressionAttributeNames["#addressLine2"] = &addressLine2String
	expressionAttributeNames["#city"] = &cityString
	expressionAttributeNames["#country"] = &countryString
	expressionAttributeNames["#details"] = &detailsString
	expressionAttributeNames["#drinksProvided"] = &drinksProvidedString
	expressionAttributeNames["#endTime"] = &endTimeString
	expressionAttributeNames["#feeForDrinks"] = &feeForDrinksString
	expressionAttributeNames["#invitesForNewInvitees"] = &invitesForNewInviteesString
	expressionAttributeNames["#latitude"] = &latitudeString
	expressionAttributeNames["#longitude"] = &longitudeString
	expressionAttributeNames["#startTime"] = &startTimeString
	expressionAttributeNames["#stateProvinceRegion"] = &stateProvinceRegionString
	expressionAttributeNames["#title"] = &titleString
	expressionAttributeNames["#zipCode"] = &zipCodeString
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var addressLine1AttributeValue = dynamodb.AttributeValue{}
	var addressLine2AttributeValue = dynamodb.AttributeValue{}
	var cityAttributeValue = dynamodb.AttributeValue{}
	var countryAttributeValue = dynamodb.AttributeValue{}
	var detailsAttributeValue = dynamodb.AttributeValue{}
	var drinksProvidedAttributeValue = dynamodb.AttributeValue{}
	var endTimeAttributeValue = dynamodb.AttributeValue{}
	var feeForDrinksAttributeValue = dynamodb.AttributeValue{}
	var invitesForNewInviteesAttributeValue = dynamodb.AttributeValue{}
	var latitudeAttributeValue = dynamodb.AttributeValue{}
	var longitudeAttributeValue = dynamodb.AttributeValue{}
	var startTimeAttributeValue = dynamodb.AttributeValue{}
	var stateProvinceRegionAttributeValue = dynamodb.AttributeValue{}
	var titleAttributeValue = dynamodb.AttributeValue{}
	var zipCodeAttributeValue = dynamodb.AttributeValue{}
	invitesForNewInviteesAttributeValue.SetN(invitesForNewInvitees)
	addressLine1AttributeValue.SetS(addressLine1)
	addressLine2AttributeValue.SetS(addressLine2)
	cityAttributeValue.SetS(city)
	countryAttributeValue.SetS(country)
	detailsAttributeValue.SetS(details)
	drinksProvidedAttributeValue.SetBOOL(drinksProvided)
	endTimeAttributeValue.SetS(endTime)
	feeForDrinksAttributeValue.SetBOOL(feeForDrinks)
	invitesForNewInviteesAttributeValue.SetN(invitesForNewInvitees)
	latitudeAttributeValue.SetN(latitude)
	longitudeAttributeValue.SetN(longitude)
	startTimeAttributeValue.SetS(startTime)
	stateProvinceRegionAttributeValue.SetS(stateProvinceRegion)
	titleAttributeValue.SetS(title)
	zipCodeAttributeValue.SetN(zipCode)
	expressionValuePlaceholders[":addressLine1"] = &addressLine1AttributeValue
	expressionValuePlaceholders[":addressLine2"] = &addressLine2AttributeValue
	expressionValuePlaceholders[":city"] = &cityAttributeValue
	expressionValuePlaceholders[":country"] = &countryAttributeValue
	expressionValuePlaceholders[":details"] = &detailsAttributeValue
	expressionValuePlaceholders[":drinksProvided"] = &drinksProvidedAttributeValue
	expressionValuePlaceholders[":endTime"] = &endTimeAttributeValue
	expressionValuePlaceholders[":feeForDrinks"] = &feeForDrinksAttributeValue
	expressionValuePlaceholders[":invitesForNewInvitees"] = &invitesForNewInviteesAttributeValue
	expressionValuePlaceholders[":latitude"] = &latitudeAttributeValue
	expressionValuePlaceholders[":longitude"] = &longitudeAttributeValue
	expressionValuePlaceholders[":startTime"] = &startTimeAttributeValue
	expressionValuePlaceholders[":stateProvinceRegion"] = &stateProvinceRegionAttributeValue
	expressionValuePlaceholders[":title"] = &titleAttributeValue
	expressionValuePlaceholders[":zipCode"] = &zipCodeAttributeValue
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #addressLine1=:addressLine1, #addressLine2=:addressLine2, #city=:city, #country=:country, #details=:details, #drinksProvided=:drinksProvided, #endTime=:endTime, #feeForDrinks=:feeForDrinks, #invitesForNewInvitees=:invitesForNewInvitees, #latitude=:latitude, #longitude=:longitude, #startTime=:startTime, #stateProvinceRegion=:stateProvinceRegion, #title=:title, #zipCode=:zipCode"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "updatePartyHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}

func deletePartyHelper(partyID string) QueryResult {
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
		queryResult.Error = "deletePartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(partyID)
	keyMap["partyID"] = &key

	var deleteItemInput = dynamodb.DeleteItemInput{}
	deleteItemInput.SetTableName("Party")
	deleteItemInput.SetKey(keyMap)

	_, err2 := getter.DynamoDB.DeleteItem(&deleteItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "deletePartyHelper function: DeleteItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}
