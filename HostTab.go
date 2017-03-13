package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// NULL : a string which contains the word "null" to replace the empty
//		string "" since dynamodb doesn't allow empty strings
const NULL = "null"

// Create a party
func createParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	name := r.Form.Get("name")
	addressLine1 := r.Form.Get("addressLine1")
	addressLine2 := r.Form.Get("addressLine2")
	city := r.Form.Get("city")
	country := r.Form.Get("country")
	details := r.Form.Get("details")
	drinksProvided, drinksProvidedConvErr := strconv.ParseBool(r.Form.Get("drinksProvided"))
	endTime := r.Form.Get("endTime")
	feeForDrinks, feeForDrinksConvErr := strconv.ParseBool(r.Form.Get("feeForDrinks"))
	invitesForNewInvitees := r.Form.Get("invitesForNewInvitees")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	partyID := strconv.FormatUint(getRandomID(), 10)
	startTime := r.Form.Get("startTime")
	stateProvinceRegion := r.Form.Get("stateProvinceRegion")
	title := r.Form.Get("title")
	zipCode := r.Form.Get("zipCode")
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
		addressLine2 = NULL
	}
	if details == "" {
		details = NULL
	}
	queryResult = createPartyHelper(facebookID, isMale, name, addressLine1, addressLine2, city, country, details, drinksProvided, endTime, feeForDrinks, invitesForNewInvitees, latitude, longitude, partyID, startTime, stateProvinceRegion, title, zipCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

// Create a party
func createBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	nameOfCreator := r.Form.Get("nameOfCreator")
	addressLine1 := r.Form.Get("addressLine1")
	addressLine2 := r.Form.Get("addressLine2")
	barID := strconv.FormatUint(getRandomID(), 10)
	city := r.Form.Get("city")
	country := r.Form.Get("country")
	details := r.Form.Get("details")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	name := r.Form.Get("name")
	phoneNumber := r.Form.Get("phoneNumber")
	stateProvinceRegion := r.Form.Get("stateProvinceRegion")
	zipCode := r.Form.Get("zipCode")
	mon := strings.Split(r.Form.Get("Mon"), ",")
	tue := strings.Split(r.Form.Get("Tue"), ",")
	wed := strings.Split(r.Form.Get("Wed"), ",")
	thu := strings.Split(r.Form.Get("Thu"), ",")
	fri := strings.Split(r.Form.Get("Fri"), ",")
	sat := strings.Split(r.Form.Get("Sat"), ",")
	sun := strings.Split(r.Form.Get("Sun"), ",")
	var scheduleForMonday = ScheduleForDay{Open: mon[0], LastCall: mon[1]}
	var scheduleForTuesday = ScheduleForDay{Open: tue[0], LastCall: tue[1]}
	var scheduleForWednesday = ScheduleForDay{Open: wed[0], LastCall: wed[1]}
	var scheduleForThursday = ScheduleForDay{Open: thu[0], LastCall: thu[1]}
	var scheduleForFriday = ScheduleForDay{Open: fri[0], LastCall: fri[1]}
	var scheduleForSaturday = ScheduleForDay{Open: sat[0], LastCall: sat[1]}
	var scheduleForSunday = ScheduleForDay{Open: sun[0], LastCall: sun[1]}
	schedule := make(map[string]ScheduleForDay)
	schedule["Monday"] = scheduleForMonday
	schedule["Tuesday"] = scheduleForTuesday
	schedule["Wednesday"] = scheduleForWednesday
	schedule["Thursday"] = scheduleForThursday
	schedule["Friday"] = scheduleForFriday
	schedule["Saturday"] = scheduleForSaturday
	schedule["Sunday"] = scheduleForSunday
	var queryResult = QueryResult{}
	if isMaleConvErr != nil {
		queryResult.Error = "createBar function: isMale parameter issue. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	if addressLine2 == "" {
		addressLine2 = NULL
	}
	if details == "" {
		details = NULL
	}
	queryResult = createBarHelper(facebookID, isMale, nameOfCreator, addressLine1, addressLine2, barID, city, country, details, latitude, longitude, name, phoneNumber, schedule, stateProvinceRegion, zipCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

// Delete a Party from the database
func deleteParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	queryResult := deletePartyHelper(partyID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

// Delete a Bar from the database
func deleteBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	queryResult := deleteBarHelper(barID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

// Update a party's information
func updateParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addressLine1 := r.Form.Get("addressLine1")
	addressLine2 := r.Form.Get("addressLine2")
	city := r.Form.Get("city")
	country := r.Form.Get("country")
	details := r.Form.Get("details")
	drinksProvided, drinksProvidedConvErr := strconv.ParseBool(r.Form.Get("drinksProvided"))
	endTime := r.Form.Get("endTime")
	feeForDrinks, feeForDrinksConvErr := strconv.ParseBool(r.Form.Get("feeForDrinks"))
	invitesForNewInvitees := r.Form.Get("invitesForNewInvitees")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	partyID := r.Form.Get("partyID")
	startTime := r.Form.Get("startTime")
	stateProvinceRegion := r.Form.Get("stateProvinceRegion")
	title := r.Form.Get("title")
	zipCode := r.Form.Get("zipCode")
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
		addressLine2 = NULL
	}
	if details == "" {
		details = NULL
	}
	queryResult = updatePartyHelper(addressLine1, addressLine2, city, country, details, drinksProvided, endTime, feeForDrinks, invitesForNewInvitees, latitude, longitude, partyID, startTime, stateProvinceRegion, title, zipCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

// Update a bar's information
func updateBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	addressLine1 := r.Form.Get("addressLine1")
	addressLine2 := r.Form.Get("addressLine2")
	barID := r.Form.Get("barID")
	city := r.Form.Get("city")
	country := r.Form.Get("country")
	details := r.Form.Get("details")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	name := r.Form.Get("name")
	phoneNumber := r.Form.Get("phoneNumber")
	stateProvinceRegion := r.Form.Get("stateProvinceRegion")
	zipCode := r.Form.Get("zipCode")
	mon := strings.Split(r.Form.Get("Mon"), ",")
	tue := strings.Split(r.Form.Get("Tue"), ",")
	wed := strings.Split(r.Form.Get("Wed"), ",")
	thu := strings.Split(r.Form.Get("Thu"), ",")
	fri := strings.Split(r.Form.Get("Fri"), ",")
	sat := strings.Split(r.Form.Get("Sat"), ",")
	sun := strings.Split(r.Form.Get("Sun"), ",")
	var scheduleForMonday = ScheduleForDay{Open: mon[0], LastCall: mon[1]}
	var scheduleForTuesday = ScheduleForDay{Open: tue[0], LastCall: tue[1]}
	var scheduleForWednesday = ScheduleForDay{Open: wed[0], LastCall: wed[1]}
	var scheduleForThursday = ScheduleForDay{Open: thu[0], LastCall: thu[1]}
	var scheduleForFriday = ScheduleForDay{Open: fri[0], LastCall: fri[1]}
	var scheduleForSaturday = ScheduleForDay{Open: sat[0], LastCall: sat[1]}
	var scheduleForSunday = ScheduleForDay{Open: sun[0], LastCall: sun[1]}
	schedule := make(map[string]ScheduleForDay)
	schedule["Monday"] = scheduleForMonday
	schedule["Tuesday"] = scheduleForTuesday
	schedule["Wednesday"] = scheduleForWednesday
	schedule["Thursday"] = scheduleForThursday
	schedule["Friday"] = scheduleForFriday
	schedule["Saturday"] = scheduleForSaturday
	schedule["Sunday"] = scheduleForSunday
	if addressLine2 == "" {
		addressLine2 = NULL
	}
	if details == "" {
		details = NULL
	}
	queryResult := updateBarHelper(addressLine1, addressLine2, barID, city, country, details, latitude, longitude, name, phoneNumber, schedule, stateProvinceRegion, zipCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

// As a host of the party, change the number of invitations
//		(doesn't have to be the same number) that the selected guests have
func setNumberOfInvitationsLeftForInvitees(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	invitees := strings.Split(r.Form.Get("invitees"), ",")
	invitationsLeft := strings.Split(r.Form.Get("invitationsLeft"), ",")
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if len(invitees) != len(invitationsLeft) {
		queryResult.Error = "Length of invitees array is not the same length as the invitationsLeft array."
	} else {
		queryResult = setNumberOfInvitationsLeftForInviteesHelper(partyID, invitees, invitationsLeft)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(queryResult)
}

func createBarHelper(facebookID string, isMale bool, nameOfCreator string, addressLine1 string, addressLine2 string, barID string, city string, country string, details string, latitude string, longitude string, name string, phoneNumber string, schedule map[string]ScheduleForDay, stateProvinceRegion string, zipCode string) QueryResult {
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
		queryResult.Error = "createBarHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionValues := make(map[string]*dynamodb.AttributeValue)
	var addressLine1AttributeValue = dynamodb.AttributeValue{}
	var addressLine2AttributeValue = dynamodb.AttributeValue{}
	var barIDAttributeValue = dynamodb.AttributeValue{}
	var cityAttributeValue = dynamodb.AttributeValue{}
	var countryAttributeValue = dynamodb.AttributeValue{}
	var detailsAttributeValue = dynamodb.AttributeValue{}
	var latitudeAttributeValue = dynamodb.AttributeValue{}
	var longitudeAttributeValue = dynamodb.AttributeValue{}
	var nameAttributeValue = dynamodb.AttributeValue{}
	var phoneNumberAttributeValue = dynamodb.AttributeValue{}
	var stateProvinceRegionAttributeValue = dynamodb.AttributeValue{}
	var zipCodeAttributeValue = dynamodb.AttributeValue{}
	addressLine1AttributeValue.SetS(addressLine1)
	addressLine2AttributeValue.SetS(addressLine2)
	barIDAttributeValue.SetN(barID)
	cityAttributeValue.SetS(city)
	countryAttributeValue.SetS(country)
	detailsAttributeValue.SetS(details)
	latitudeAttributeValue.SetN(latitude)
	longitudeAttributeValue.SetN(longitude)
	nameAttributeValue.SetS(name)
	phoneNumberAttributeValue.SetS(phoneNumber)
	stateProvinceRegionAttributeValue.SetS(stateProvinceRegion)
	zipCodeAttributeValue.SetN(zipCode)
	expressionValues["addressLine1"] = &addressLine1AttributeValue
	expressionValues["addressLine2"] = &addressLine2AttributeValue
	expressionValues["barID"] = &barIDAttributeValue
	expressionValues["city"] = &cityAttributeValue
	expressionValues["country"] = &countryAttributeValue
	expressionValues["details"] = &detailsAttributeValue
	expressionValues["latitude"] = &latitudeAttributeValue
	expressionValues["longitude"] = &longitudeAttributeValue
	expressionValues["name"] = &nameAttributeValue
	expressionValues["phoneNumber"] = &phoneNumberAttributeValue
	expressionValues["stateProvinceRegion"] = &stateProvinceRegionAttributeValue
	expressionValues["zipCode"] = &zipCodeAttributeValue

	// set yourself as an attendee to your own bar so that you can rate it
	attendeesMap := make(map[string]*dynamodb.AttributeValue)
	var attendees = dynamodb.AttributeValue{}
	attendeeMap := make(map[string]*dynamodb.AttributeValue)
	var attendee = dynamodb.AttributeValue{}
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameOfCreatorAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	isMaleAttribute.SetBOOL(isMale)
	nameOfCreatorAttribute.SetS(nameOfCreator)
	ratingAttribute.SetS("N")
	statusAttribute.SetS("G")
	timeLastRatedAttribute.SetS("01/01/2000 00:00:00")
	attendeeMap["isMale"] = &isMaleAttribute
	attendeeMap["name"] = &nameOfCreatorAttribute
	attendeeMap["rating"] = &ratingAttribute
	attendeeMap["status"] = &statusAttribute
	attendeeMap["timeLastRated"] = &timeLastRatedAttribute
	attendee.SetM(attendeeMap)
	attendeesMap[facebookID] = &attendee
	attendees.SetM(attendeesMap)
	expressionValues["attendees"] = &attendees

	scheduleMap := make(map[string]*dynamodb.AttributeValue)
	var theSchedule = dynamodb.AttributeValue{}
	scheduleForMonday := make(map[string]*dynamodb.AttributeValue)
	scheduleForTuesday := make(map[string]*dynamodb.AttributeValue)
	scheduleForWednesday := make(map[string]*dynamodb.AttributeValue)
	scheduleForThursday := make(map[string]*dynamodb.AttributeValue)
	scheduleForFriday := make(map[string]*dynamodb.AttributeValue)
	scheduleForSaturday := make(map[string]*dynamodb.AttributeValue)
	scheduleForSunday := make(map[string]*dynamodb.AttributeValue)
	var openMonday = dynamodb.AttributeValue{}
	var lastCallMonday = dynamodb.AttributeValue{}
	var openTuesday = dynamodb.AttributeValue{}
	var lastCallTuesday = dynamodb.AttributeValue{}
	var openWednesday = dynamodb.AttributeValue{}
	var lastCallWednesday = dynamodb.AttributeValue{}
	var openThursday = dynamodb.AttributeValue{}
	var lastCallThursday = dynamodb.AttributeValue{}
	var openFriday = dynamodb.AttributeValue{}
	var lastCallFriday = dynamodb.AttributeValue{}
	var openSaturday = dynamodb.AttributeValue{}
	var lastCallSaturday = dynamodb.AttributeValue{}
	var openSunday = dynamodb.AttributeValue{}
	var lastCallSunday = dynamodb.AttributeValue{}
	openMonday.SetS(schedule["Monday"].Open)
	lastCallMonday.SetS(schedule["Monday"].LastCall)
	openTuesday.SetS(schedule["Tuesday"].Open)
	lastCallTuesday.SetS(schedule["Tuesday"].LastCall)
	openWednesday.SetS(schedule["Wednesday"].Open)
	lastCallWednesday.SetS(schedule["Wednesday"].LastCall)
	openThursday.SetS(schedule["Thursday"].Open)
	lastCallThursday.SetS(schedule["Thursday"].LastCall)
	openFriday.SetS(schedule["Friday"].Open)
	lastCallFriday.SetS(schedule["Friday"].LastCall)
	openSaturday.SetS(schedule["Saturday"].Open)
	lastCallSaturday.SetS(schedule["Saturday"].LastCall)
	openSunday.SetS(schedule["Sunday"].Open)
	lastCallSunday.SetS(schedule["Sunday"].LastCall)
	scheduleForMonday["open"] = &openMonday
	scheduleForMonday["lastCall"] = &lastCallMonday
	scheduleForTuesday["open"] = &openTuesday
	scheduleForTuesday["lastCall"] = &lastCallTuesday
	scheduleForWednesday["open"] = &openWednesday
	scheduleForWednesday["lastCall"] = &lastCallWednesday
	scheduleForThursday["open"] = &openThursday
	scheduleForThursday["lastCall"] = &lastCallThursday
	scheduleForFriday["open"] = &openFriday
	scheduleForFriday["lastCall"] = &lastCallFriday
	scheduleForSaturday["open"] = &openSaturday
	scheduleForSaturday["lastCall"] = &lastCallSaturday
	scheduleForSunday["open"] = &openSunday
	scheduleForSunday["lastCall"] = &lastCallSunday
	var scheduleForMondayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForTuesdayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForWednesdayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForThursdayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForFridayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForSaturdayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForSundayAttributeValue = dynamodb.AttributeValue{}
	scheduleForMondayAttributeValue.SetM(scheduleForMonday)
	scheduleForTuesdayAttributeValue.SetM(scheduleForTuesday)
	scheduleForWednesdayAttributeValue.SetM(scheduleForWednesday)
	scheduleForThursdayAttributeValue.SetM(scheduleForThursday)
	scheduleForFridayAttributeValue.SetM(scheduleForFriday)
	scheduleForSaturdayAttributeValue.SetM(scheduleForSaturday)
	scheduleForSundayAttributeValue.SetM(scheduleForSunday)
	scheduleMap["Monday"] = &scheduleForMondayAttributeValue
	scheduleMap["Tuesday"] = &scheduleForTuesdayAttributeValue
	scheduleMap["Wednesday"] = &scheduleForWednesdayAttributeValue
	scheduleMap["Thursday"] = &scheduleForThursdayAttributeValue
	scheduleMap["Friday"] = &scheduleForFridayAttributeValue
	scheduleMap["Saturday"] = &scheduleForSaturdayAttributeValue
	scheduleMap["Sunday"] = &scheduleForSundayAttributeValue
	theSchedule.SetM(scheduleMap)
	expressionValues["schedule"] = &theSchedule

	hostsMap := make(map[string]*dynamodb.AttributeValue)
	var hosts = dynamodb.AttributeValue{}
	hostMap := make(map[string]*dynamodb.AttributeValue)
	var host = dynamodb.AttributeValue{}
	var isMainHostAttribute = dynamodb.AttributeValue{}
	isMainHostAttribute.SetBOOL(true)
	nameOfCreatorAttribute.SetS(nameOfCreator)
	hostMap["isMainHost"] = &isMainHostAttribute
	hostMap["name"] = &nameOfCreatorAttribute
	host.SetM(hostMap)
	hostsMap[facebookID] = &host
	hosts.SetM(hostsMap)
	expressionValues["hosts"] = &hosts

	var putItemInput = dynamodb.PutItemInput{}
	putItemInput.SetTableName("Bar")
	putItemInput.SetItem(expressionValues)
	_, err2 := getter.DynamoDB.PutItem(&putItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "createBarHelper function: PutItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall

	// Now we need to update the person's information to let them
	//     know that they are hosting this bar.
	expressionAttributeNames := make(map[string]*string)
	barHostFor := "barHostFor"
	expressionAttributeNames["#barHostFor"] = &barHostFor
	expressionAttributeNames["#barID"] = &barID
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var barIDBoolAttribute = dynamodb.AttributeValue{}
	barIDBoolAttribute.SetBOOL(true)
	expressionValuePlaceholders[":bool"] = &barIDBoolAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "SET #barHostFor.#barID=:bool"
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall2 = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall2.Error = "createBarHelper function: UpdateItem error (probable cause: your facebookID isn't in the database). " + updateItemOutputErr.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
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
		queryResult.Error = "createPartyHelper function: session creation error. " + err.Error()
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
	queryResult.DynamodbCalls = nil
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
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func updateBarHelper(addressLine1 string, addressLine2 string, barID string, city string, country string, details string, latitude string, longitude string, name string, phoneNumber string, schedule map[string]ScheduleForDay, stateProvinceRegion string, zipCode string) QueryResult {
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
		queryResult.Error = "updateBarHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally

	// switch: startTime to openAt
	expressionAttributeNames := make(map[string]*string)
	var addressLine1String = "addressLine1"
	var addressLine2String = "addressLine2"
	var cityString = "city"
	var countryString = "country"
	var detailsString = "details"
	var latitudeString = "latitude"
	var longitudeString = "longitude"
	var nameString = "name"
	var phoneNumberString = "phoneNumber"
	var scheduleString = "schedule"
	var stateProvinceRegionString = "stateProvinceRegion"
	var zipCodeString = "zipCode"
	expressionAttributeNames["#addressLine1"] = &addressLine1String
	expressionAttributeNames["#addressLine2"] = &addressLine2String
	expressionAttributeNames["#city"] = &cityString
	expressionAttributeNames["#country"] = &countryString
	expressionAttributeNames["#details"] = &detailsString
	expressionAttributeNames["#latitude"] = &latitudeString
	expressionAttributeNames["#longitude"] = &longitudeString
	expressionAttributeNames["#name"] = &nameString
	expressionAttributeNames["#phoneNumber"] = &phoneNumberString
	expressionAttributeNames["#schedule"] = &scheduleString
	expressionAttributeNames["#stateProvinceRegion"] = &stateProvinceRegionString
	expressionAttributeNames["#zipCode"] = &zipCodeString
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var addressLine1AttributeValue = dynamodb.AttributeValue{}
	var addressLine2AttributeValue = dynamodb.AttributeValue{}
	var cityAttributeValue = dynamodb.AttributeValue{}
	var countryAttributeValue = dynamodb.AttributeValue{}
	var detailsAttributeValue = dynamodb.AttributeValue{}
	var latitudeAttributeValue = dynamodb.AttributeValue{}
	var longitudeAttributeValue = dynamodb.AttributeValue{}
	var nameAttributeValue = dynamodb.AttributeValue{}
	var phoneNumberAttributeValue = dynamodb.AttributeValue{}
	var stateProvinceRegionAttributeValue = dynamodb.AttributeValue{}
	var zipCodeAttributeValue = dynamodb.AttributeValue{}
	addressLine1AttributeValue.SetS(addressLine1)
	addressLine2AttributeValue.SetS(addressLine2)
	cityAttributeValue.SetS(city)
	countryAttributeValue.SetS(country)
	detailsAttributeValue.SetS(details)
	latitudeAttributeValue.SetN(latitude)
	longitudeAttributeValue.SetN(longitude)
	nameAttributeValue.SetS(name)
	phoneNumberAttributeValue.SetS(phoneNumber)
	stateProvinceRegionAttributeValue.SetS(stateProvinceRegion)
	zipCodeAttributeValue.SetN(zipCode)
	expressionValuePlaceholders[":addressLine1"] = &addressLine1AttributeValue
	expressionValuePlaceholders[":addressLine2"] = &addressLine2AttributeValue
	expressionValuePlaceholders[":city"] = &cityAttributeValue
	expressionValuePlaceholders[":country"] = &countryAttributeValue
	expressionValuePlaceholders[":details"] = &detailsAttributeValue
	expressionValuePlaceholders[":latitude"] = &latitudeAttributeValue
	expressionValuePlaceholders[":longitude"] = &longitudeAttributeValue
	expressionValuePlaceholders[":name"] = &nameAttributeValue
	expressionValuePlaceholders[":phoneNumber"] = &phoneNumberAttributeValue
	expressionValuePlaceholders[":stateProvinceRegion"] = &stateProvinceRegionAttributeValue
	expressionValuePlaceholders[":zipCode"] = &zipCodeAttributeValue

	scheduleMap := make(map[string]*dynamodb.AttributeValue)
	var theSchedule = dynamodb.AttributeValue{}
	scheduleForMonday := make(map[string]*dynamodb.AttributeValue)
	scheduleForTuesday := make(map[string]*dynamodb.AttributeValue)
	scheduleForWednesday := make(map[string]*dynamodb.AttributeValue)
	scheduleForThursday := make(map[string]*dynamodb.AttributeValue)
	scheduleForFriday := make(map[string]*dynamodb.AttributeValue)
	scheduleForSaturday := make(map[string]*dynamodb.AttributeValue)
	scheduleForSunday := make(map[string]*dynamodb.AttributeValue)
	var openMonday = dynamodb.AttributeValue{}
	var lastCallMonday = dynamodb.AttributeValue{}
	var openTuesday = dynamodb.AttributeValue{}
	var lastCallTuesday = dynamodb.AttributeValue{}
	var openWednesday = dynamodb.AttributeValue{}
	var lastCallWednesday = dynamodb.AttributeValue{}
	var openThursday = dynamodb.AttributeValue{}
	var lastCallThursday = dynamodb.AttributeValue{}
	var openFriday = dynamodb.AttributeValue{}
	var lastCallFriday = dynamodb.AttributeValue{}
	var openSaturday = dynamodb.AttributeValue{}
	var lastCallSaturday = dynamodb.AttributeValue{}
	var openSunday = dynamodb.AttributeValue{}
	var lastCallSunday = dynamodb.AttributeValue{}
	openMonday.SetS(schedule["Monday"].Open)
	lastCallMonday.SetS(schedule["Monday"].LastCall)
	openTuesday.SetS(schedule["Tuesday"].Open)
	lastCallTuesday.SetS(schedule["Tuesday"].LastCall)
	openWednesday.SetS(schedule["Wednesday"].Open)
	lastCallWednesday.SetS(schedule["Wednesday"].LastCall)
	openThursday.SetS(schedule["Thursday"].Open)
	lastCallThursday.SetS(schedule["Thursday"].LastCall)
	openFriday.SetS(schedule["Friday"].Open)
	lastCallFriday.SetS(schedule["Friday"].LastCall)
	openSaturday.SetS(schedule["Saturday"].Open)
	lastCallSaturday.SetS(schedule["Saturday"].LastCall)
	openSunday.SetS(schedule["Sunday"].Open)
	lastCallSunday.SetS(schedule["Sunday"].LastCall)
	scheduleForMonday["open"] = &openMonday
	scheduleForMonday["lastCall"] = &lastCallMonday
	scheduleForTuesday["open"] = &openTuesday
	scheduleForTuesday["lastCall"] = &lastCallTuesday
	scheduleForWednesday["open"] = &openWednesday
	scheduleForWednesday["lastCall"] = &lastCallWednesday
	scheduleForThursday["open"] = &openThursday
	scheduleForThursday["lastCall"] = &lastCallThursday
	scheduleForFriday["open"] = &openFriday
	scheduleForFriday["lastCall"] = &lastCallFriday
	scheduleForSaturday["open"] = &openSaturday
	scheduleForSaturday["lastCall"] = &lastCallSaturday
	scheduleForSunday["open"] = &openSunday
	scheduleForSunday["lastCall"] = &lastCallSunday
	var scheduleForMondayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForTuesdayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForWednesdayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForThursdayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForFridayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForSaturdayAttributeValue = dynamodb.AttributeValue{}
	var scheduleForSundayAttributeValue = dynamodb.AttributeValue{}
	scheduleForMondayAttributeValue.SetM(scheduleForMonday)
	scheduleForTuesdayAttributeValue.SetM(scheduleForTuesday)
	scheduleForWednesdayAttributeValue.SetM(scheduleForWednesday)
	scheduleForThursdayAttributeValue.SetM(scheduleForThursday)
	scheduleForFridayAttributeValue.SetM(scheduleForFriday)
	scheduleForSaturdayAttributeValue.SetM(scheduleForSaturday)
	scheduleForSundayAttributeValue.SetM(scheduleForSunday)
	scheduleMap["Monday"] = &scheduleForMondayAttributeValue
	scheduleMap["Tuesday"] = &scheduleForTuesdayAttributeValue
	scheduleMap["Wednesday"] = &scheduleForWednesdayAttributeValue
	scheduleMap["Thursday"] = &scheduleForThursdayAttributeValue
	scheduleMap["Friday"] = &scheduleForFridayAttributeValue
	scheduleMap["Saturday"] = &scheduleForSaturdayAttributeValue
	scheduleMap["Sunday"] = &scheduleForSundayAttributeValue
	theSchedule.SetM(scheduleMap)
	expressionValuePlaceholders[":schedule"] = &theSchedule

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #addressLine1=:addressLine1, #addressLine2=:addressLine2, #city=:city, #country=:country, #details=:details, #latitude=:latitude, #longitude=:longitude, #name=:name, #phoneNumber=:phoneNumber, #schedule=:schedule, #stateProvinceRegion=:stateProvinceRegion, #zipCode=:zipCode"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "updateBarHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func deletePartyHelper(partyID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 1)

	getPartyQueryResult := getParty(partyID)
	if getPartyQueryResult.Succeeded == false {
		return getPartyQueryResult
	}
	queryResult1 := removePartyFromInvitedToMapInPersonTableForEveryoneInvitedToParty(partyID, getPartyQueryResult.Parties[0].Invitees)
	if queryResult1.Succeeded == false {
		return queryResult1
	}
	queryResult2 := removePartyFromPartyHostForMapInPersonTableForEveryHostOfTheParty(partyID, getPartyQueryResult.Parties[0].Hosts)
	if queryResult2.Succeeded == false {
		return queryResult2
	}
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
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func removePartyFromInvitedToMapInPersonTableForEveryoneInvitedToParty(partyID string, invitees map[string]*Invitee) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, len(invitees))
	i := 0
	for k := range invitees {
		removeInvitationQueryResult := removePartyInviteFromInvitedToMapInPersonTable(k, partyID)
		if removeInvitationQueryResult.Succeeded == false {
			queryResult.DynamodbCalls[i] = removeInvitationQueryResult.DynamodbCalls[0]
			queryResult.Error = removeInvitationQueryResult.Error
			return queryResult
		}
		queryResult.DynamodbCalls[i] = removeInvitationQueryResult.DynamodbCalls[0]
		i++
	}
	queryResult.Succeeded = true
	return queryResult
}

func removePartyFromPartyHostForMapInPersonTableForEveryHostOfTheParty(partyID string, hosts map[string]*Host) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, len(hosts))
	i := 0
	for k := range hosts {
		qr1 := removePartyFromPartyHostForMapInPersonTable(k, partyID)
		if qr1.Succeeded == false {
			queryResult.DynamodbCalls[i] = qr1.DynamodbCalls[0]
			queryResult.Error = qr1.Error
			return queryResult
		}
		queryResult.DynamodbCalls[i] = qr1.DynamodbCalls[0]
		i++
	}
	queryResult.Succeeded = true
	return queryResult
}

func removePartyInviteFromInvitedToMapInPersonTable(facebookID string, partyID string) QueryResult {
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
		queryResult.Error = "removeInvitationToPartyInPersonTable function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var invitedTo = "invitedTo"
	expressionAttributeNames["#invitedTo"] = &invitedTo
	expressionAttributeNames["#partyID"] = &partyID

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "Remove #invitedTo.#partyID"
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall.Error = "removeInvitationToPartyInPersonTable function: UpdateItem error (probable cause: this facebookID isn't in the database). " + updateItemOutputErr.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}

func removePartyFromPartyHostForMapInPersonTable(facebookID string, partyID string) QueryResult {
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
		queryResult.Error = "removePartyFromPartyHostForMapInPersonTable function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var partyHostFor = "partyHostFor"
	expressionAttributeNames["#partyHostFor"] = &partyHostFor
	expressionAttributeNames["#partyID"] = &partyID

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "Remove #partyHostFor.#partyID"
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall.Error = "removePartyFromPartyHostForMapInPersonTable function: UpdateItem error (probable cause: this facebookID isn't in the database). " + updateItemOutputErr.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}

func deleteBarHelper(barID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, 1)

	getBarQueryResult := getBar(barID)
	if getBarQueryResult.Succeeded == false {
		return getBarQueryResult
	}
	removeHostsQueryResult := removeBarFromBarHostForMapInPersonTableForEveryHostOfTheBar(barID, getBarQueryResult.Bars[0].Hosts)
	if removeHostsQueryResult.Succeeded == false {
		return removeHostsQueryResult
	}
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		queryResult.Error = "deleteBarHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(barID)
	keyMap["barID"] = &key

	var deleteItemInput = dynamodb.DeleteItemInput{}
	deleteItemInput.SetTableName("Bar")
	deleteItemInput.SetKey(keyMap)

	_, err2 := getter.DynamoDB.DeleteItem(&deleteItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "deleteBarHelper function: DeleteItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func removeBarFromBarHostForMapInPersonTableForEveryHostOfTheBar(barID string, hosts map[string]*Host) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	queryResult.DynamodbCalls = make([]DynamodbCall, len(hosts))
	i := 0
	for k := range hosts {
		qr1 := removeBarFromBarHostForMapInPersonTable(k, barID)
		if qr1.Succeeded == false {
			queryResult.DynamodbCalls[i] = qr1.DynamodbCalls[0]
			queryResult.Error = qr1.Error
			return queryResult
		}
		queryResult.DynamodbCalls[i] = qr1.DynamodbCalls[0]
		i++
	}
	queryResult.Succeeded = true
	return queryResult
}

func removeBarFromBarHostForMapInPersonTable(facebookID string, barID string) QueryResult {
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
		queryResult.Error = "removeBarFromBarHostForMapInPersonTable function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var barHostFor = "barHostFor"
	expressionAttributeNames["#barHostFor"] = &barHostFor
	expressionAttributeNames["#barID"] = &barID

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(facebookID)
	keyMap["facebookID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Person")
	updateExpression := "Remove #barHostFor.#barID"
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall.Error = "removeBarFromBarHostForMapInPersonTable function: UpdateItem error (probable cause: this facebookID isn't in the database). " + updateItemOutputErr.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	queryResult.Succeeded = true
	return queryResult
}

func getParty(partyID string) QueryResult {
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
		queryResult.Error = "getParty function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(partyID)
	keyMap["partyID"] = &key
	getItemInput.SetKey(keyMap)
	getItemInput.SetTableName("Party")
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getParty function: GetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := getItemOutput.Item
	var party = PartyData{}
	jsonErr := dynamodbattribute.UnmarshalMap(data, &party)
	if jsonErr != nil {
		queryResult.Error = "getParty function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.Parties = make([]PartyData, 1)
	queryResult.Parties[0] = party
	queryResult.Succeeded = true
	return queryResult
}

func getBar(barID string) QueryResult {
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
		queryResult.Error = "getBar function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(barID)
	keyMap["barID"] = &key
	getItemInput.SetKey(keyMap)
	getItemInput.SetTableName("Bar")
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getBar function: GetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := getItemOutput.Item
	var bar = BarData{}
	jsonErr := dynamodbattribute.UnmarshalMap(data, &bar)
	if jsonErr != nil {
		queryResult.Error = "getBar function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.Bars = make([]BarData, 1)
	queryResult.Bars[0] = bar
	queryResult.Succeeded = true
	return queryResult
}

func setNumberOfInvitationsLeftForInviteesHelper(partyID string, invitees []string, invitationsLeft []string) QueryResult {
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
		queryResult.Error = "setNumberOfInvitationsLeftForInviteesHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var inviteesString = "invitees"
	var numberOfInvitationsLeftString = "numberOfInvitationsLeft"
	expressionAttributeNames["#invitees"] = &inviteesString
	expressionAttributeNames["#numberOfInvitationsLeft"] = &numberOfInvitationsLeftString
	for i := 0; i < len(invitees); i++ {
		expressionAttributeNames["#"+invitees[i]] = &invitees[i]
	}

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	for i := 0; i < len(invitationsLeft); i++ {
		var numberOfInvitationsLeftAtributeValue = dynamodb.AttributeValue{}
		numberOfInvitationsLeftAtributeValue.SetN(invitationsLeft[i])
		expressionValuePlaceholders[":"+invitees[i]] = &numberOfInvitationsLeftAtributeValue
	}

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetN(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	// "SET #invitees.#facebookID1.#numberOfInvitationsLeft=:numberOfInvitationsLeft1, #invitees.#facebookID2.#numberOfInvitationsLeft=:numberOfInvitationsLeft2 ...."
	updateExpression := "SET #invitees.#" + invitees[0] + ".#numberOfInvitationsLeft" + "=:" + invitees[0]
	for i := 1; i < len(invitationsLeft); i++ {
		updateExpression = updateExpression + ", " + "#invitees.#" + invitees[i] + ".#numberOfInvitationsLeft" + "=:" + invitees[i]
	}
	fmt.Println(updateExpression)
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall.Error = "setNumberOfInvitationsLeftForInviteesHelper function: UpdateItem error (probable cause: this partyID isn't in the database). " + updateItemOutputErr.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}
