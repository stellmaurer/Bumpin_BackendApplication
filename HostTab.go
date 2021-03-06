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

func getClaimKey(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	queryResult := getClaimKeyHelper(key)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func getBarKey(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	queryResult := getBarKeyHelper(key)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func deleteBarKey(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	queryResult := deleteBarKeyHelper(key)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Create a party
func createParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	name := r.Form.Get("name")
	address := r.Form.Get("address")
	details := r.Form.Get("details")
	drinksProvided, drinksProvidedConvErr := strconv.ParseBool(r.Form.Get("drinksProvided"))
	endTime := r.Form.Get("endTime")
	feeForDrinks, feeForDrinksConvErr := strconv.ParseBool(r.Form.Get("feeForDrinks"))
	invitesForNewInvitees := r.Form.Get("invitesForNewInvitees")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	partyID := strconv.FormatUint(getRandomID(), 10)
	startTime := r.Form.Get("startTime")
	title := r.Form.Get("title")
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
	if details == "" {
		details = NULL
	}
	queryResult = createPartyHelper(facebookID, isMale, name, address, details, drinksProvided, endTime, feeForDrinks, invitesForNewInvitees, latitude, longitude, partyID, startTime, title)

	if queryResult.Succeeded == true {
		var addHostsQueryResult = askFriendsToHostPartyWithYou(r, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, addHostsQueryResult)

		var updateInvitationsQueryResult = updateInvitationsListAsHostForParty(r, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, updateInvitationsQueryResult)
	}
	if queryResult.Succeeded == true {
		queryResult.Error = partyID
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Create a bar
func createBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barKey := r.Form.Get("barKey")
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	nameOfCreator := r.Form.Get("nameOfCreator")
	address := r.Form.Get("address")
	attendeesMapCleanUpHourInZulu := r.Form.Get("attendeesMapCleanUpHourInZulu")
	barID := strconv.FormatUint(getRandomID(), 10)
	details := r.Form.Get("details")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	name := r.Form.Get("name")
	phoneNumber := r.Form.Get("phoneNumber")
	timeZone := r.Form.Get("timeZone")
	googlePlaceID := r.Form.Get("googlePlaceID")
	if googlePlaceID == "" {
		googlePlaceID = "-1"
	}

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
	if details == "" {
		details = NULL
	}
	queryResult = createBarHelper(barKey, facebookID, isMale, nameOfCreator, address, attendeesMapCleanUpHourInZulu, barID, details, latitude, longitude, name, phoneNumber, schedule, timeZone, googlePlaceID)

	if queryResult.Succeeded == true {
		var addHostsQueryResult = askFriendsToHostBarWithYou(r, barID)
		queryResult = convertTwoQueryResultsToOne(queryResult, addHostsQueryResult)
	}
	if queryResult.Succeeded == true {
		queryResult.Error = barID
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Delete a Party from the database
func deleteParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	queryResult := deletePartyHelper(partyID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Delete a Bar from the database
func deleteBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	queryResult := deleteBarHelper(barID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Update a party's information
func updateParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	address := r.Form.Get("address")
	details := r.Form.Get("details")
	drinksProvided, drinksProvidedConvErr := strconv.ParseBool(r.Form.Get("drinksProvided"))
	endTime := r.Form.Get("endTime")
	feeForDrinks, feeForDrinksConvErr := strconv.ParseBool(r.Form.Get("feeForDrinks"))
	invitesForNewInvitees := r.Form.Get("invitesForNewInvitees")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	partyID := r.Form.Get("partyID")
	startTime := r.Form.Get("startTime")
	title := r.Form.Get("title")
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
	if details == "" {
		details = NULL
	}
	queryResult = updatePartyHelper(address, details, drinksProvided, endTime, feeForDrinks, invitesForNewInvitees, latitude, longitude, partyID, startTime, title)

	if queryResult.Succeeded == true {
		var updateHostsQueryResult = updateHostListForParty(r, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, updateHostsQueryResult)

		var updateInvitationsQueryResult = updateInvitationsListAsHostForParty(r, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, updateInvitationsQueryResult)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Update a bar's information
func updateBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	address := r.Form.Get("address")
	attendeesMapCleanUpHourInZulu := r.Form.Get("attendeesMapCleanUpHourInZulu")
	barID := r.Form.Get("barID")
	details := r.Form.Get("details")
	latitude := r.Form.Get("latitude")
	longitude := r.Form.Get("longitude")
	name := r.Form.Get("name")
	phoneNumber := r.Form.Get("phoneNumber")
	timeZone := r.Form.Get("timeZone")
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
	if details == "" {
		details = NULL
	}
	queryResult := updateBarHelper(address, attendeesMapCleanUpHourInZulu, barID, details, latitude, longitude, name, phoneNumber, schedule, timeZone)

	if queryResult.Succeeded == true {
		var updateHostsQueryResult = updateHostListForBar(r, barID)
		queryResult = convertTwoQueryResultsToOne(queryResult, updateHostsQueryResult)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Ask friends to host the Party with you
func askFriendsToHostPartyWithYou(r *http.Request, partyID string) QueryResult {
	r.ParseForm()

	var queryResult = QueryResult{}
	queryResult.Succeeded = true

	var hostListFacebookIDs []string
	var hostListNames []string

	if r.Form.Get("hostListFacebookIDs") != "" {
		hostListFacebookIDs = strings.Split(r.Form.Get("hostListFacebookIDs"), ",")
	}
	if r.Form.Get("hostListNames") != "" {
		hostListNames = strings.Split(r.Form.Get("hostListNames"), ",")
	}
	if len(hostListFacebookIDs) != len(hostListNames) {
		queryResult.Error = "askFriendsToHostPartyWithYou function: HTTP post request parameter issues (hosts lists): facebookID array and name array aren't the same length."
		queryResult.Succeeded = false
		return queryResult
	}

	for i := 0; i < len(hostListFacebookIDs); i++ {
		askFriendQueryResult := askFriendToHostPartyWithYouHelper(partyID, hostListFacebookIDs[i], hostListNames[i])
		queryResult = convertTwoQueryResultsToOne(queryResult, askFriendQueryResult)
	}

	if queryResult.Succeeded == true {
		hostIsMale, _ := strconv.ParseBool(r.Form.Get("isMale"))
		genderString := "her"
		if hostIsMale == true {
			genderString = "him"
		}
		message := r.Form.Get("name") + " wants you to host a party with " + genderString + "."
		sendPushNotificationsQueryResult := createAndSendNotificationsToThesePeople(hostListFacebookIDs, message, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, sendPushNotificationsQueryResult)
	}
	return queryResult
}

// Ask a friend to host the Party with you
func askFriendsToHostBarWithYou(r *http.Request, barID string) QueryResult {
	r.ParseForm()

	var queryResult = QueryResult{}
	queryResult.Succeeded = true

	var hostListFacebookIDs []string
	var hostListNames []string

	if r.Form.Get("hostListFacebookIDs") != "" {
		hostListFacebookIDs = strings.Split(r.Form.Get("hostListFacebookIDs"), ",")
	}
	if r.Form.Get("hostListNames") != "" {
		hostListNames = strings.Split(r.Form.Get("hostListNames"), ",")
	}
	if len(hostListFacebookIDs) != len(hostListNames) {
		queryResult.Error = "askFriendsToHostBarWithYou function: HTTP post request parameter issues (hosts lists): facebookID array and name array aren't the same length."
		queryResult.Succeeded = false
		return queryResult
	}

	for i := 0; i < len(hostListFacebookIDs); i++ {
		askFriendQueryResult := askFriendToHostBarWithYouHelper(barID, hostListFacebookIDs[i], hostListNames[i])
		queryResult = convertTwoQueryResultsToOne(queryResult, askFriendQueryResult)
	}

	if queryResult.Succeeded == true {
		hostIsMale, _ := strconv.ParseBool(r.Form.Get("isMale"))
		genderString := "her"
		if hostIsMale == true {
			genderString = "him"
		}
		message := r.Form.Get("nameOfCreator") + " wants you to host a bar with " + genderString + "."
		sendPushNotificationsQueryResult := createAndSendNotificationsToThesePeople(hostListFacebookIDs, message, barID)
		queryResult = convertTwoQueryResultsToOne(queryResult, sendPushNotificationsQueryResult)
	}
	return queryResult
}

// Remove a friend from hosting the Party with you
func removePartyHost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	facebookID := r.Form.Get("facebookID")
	queryResult := removePartyHostHelper(partyID, facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Remove a host of a Bar
func removeBarHost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	facebookID := r.Form.Get("facebookID")
	queryResult := removeBarHostHelper(barID, facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Accept your friend's invitation to become a host of the Party
func acceptInvitationToHostParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	facebookID := r.Form.Get("facebookID")
	isMale, isMaleConvErr := strconv.ParseBool(r.Form.Get("isMale"))
	name := r.Form.Get("name")
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	if isMaleConvErr != nil {
		queryResult.Error = "acceptInvitationToHostParty function: HTTP post request isMale parameter issue. " + isMaleConvErr.Error()
		json.NewEncoder(w).Encode(queryResult)
		return
	}
	queryResult = acceptInvitationToHostPartyHelper(partyID, facebookID, isMale, name)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Accept your friend's invitation to become a host of the Bar
func acceptInvitationToHostBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	facebookID := r.Form.Get("facebookID")
	queryResult := acceptInvitationToHostBarHelper(barID, facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Decline your friend's invitation to become a host of the Party
func declineInvitationToHostParty(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	partyID := r.Form.Get("partyID")
	facebookID := r.Form.Get("facebookID")
	queryResult := declineInvitationToHostPartyHelper(partyID, facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

// Decline your friend's invitation to become a host of the Bar
func declineInvitationToHostBar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	barID := r.Form.Get("barID")
	facebookID := r.Form.Get("facebookID")
	queryResult := declineInvitationToHostBarHelper(barID, facebookID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func updateInvitationsListAsHostForParty(r *http.Request, partyID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	r.ParseForm()
	var numberOfInvitesToGive = r.Form.Get("invitesForNewInvitees")
	var additionsListFacebookID []string
	var additionsListIsMaleString []string
	var additionsListName []string
	var removalsListFacebookID []string

	if r.Form.Get("additionsListFacebookID") != "" {
		additionsListFacebookID = strings.Split(r.Form.Get("additionsListFacebookID"), ",")
	}
	if r.Form.Get("additionsListIsMale") != "" {
		additionsListIsMaleString = strings.Split(r.Form.Get("additionsListIsMale"), ",")
	}
	// convert IsMale string array to IsMale bool array
	var additionsListIsMale = make([]bool, len(additionsListIsMaleString))
	for i := 0; i < len(additionsListIsMaleString); i++ {
		isMale, isMaleConvErr := strconv.ParseBool(additionsListIsMaleString[i])
		if isMaleConvErr != nil {
			queryResult.Error = "updateInvitationsListAsHostForParty function: HTTP post request isMale parameter issue. " + isMaleConvErr.Error()
			return queryResult
		}
		additionsListIsMale[i] = isMale
	}
	if r.Form.Get("additionsListName") != "" {
		additionsListName = strings.Split(r.Form.Get("additionsListName"), ",")
	}
	if r.Form.Get("removalsListFacebookID") != "" {
		removalsListFacebookID = strings.Split(r.Form.Get("removalsListFacebookID"), ",")
	}
	if (len(additionsListFacebookID) != len(additionsListIsMale)) || (len(additionsListIsMale) != len(additionsListName)) {
		queryResult.Error = "updateInvitationsListAsHostForParty function: HTTP post request parameter issues (additions lists): facebookID array, isMale array, and name array aren't the same length."
		return queryResult
	}
	queryResult = updateInvitationsListAsHostForPartyHelper(partyID, numberOfInvitesToGive, additionsListFacebookID, additionsListIsMale, additionsListName, removalsListFacebookID)

	if queryResult.Succeeded == true {
		message := r.Form.Get("name") + " invited you to a party."
		sendPushNotificationsQueryResult := createAndSendNotificationsToThesePeople(additionsListFacebookID, message, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, sendPushNotificationsQueryResult)
	}
	return queryResult
}

func updateHostListForParty(r *http.Request, partyID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	r.ParseForm()
	var hostsToAddFacebookIDs []string
	var hostsToAddNames []string
	var hostsToRemoveFacebookIDs []string

	if r.Form.Get("hostsToAddFacebookIDs") != "" {
		hostsToAddFacebookIDs = strings.Split(r.Form.Get("hostsToAddFacebookIDs"), ",")
	}
	if r.Form.Get("hostsToAddNames") != "" {
		hostsToAddNames = strings.Split(r.Form.Get("hostsToAddNames"), ",")
	}
	if r.Form.Get("hostsToRemoveFacebookIDs") != "" {
		hostsToRemoveFacebookIDs = strings.Split(r.Form.Get("hostsToRemoveFacebookIDs"), ",")
	}
	if len(hostsToAddFacebookIDs) != len(hostsToAddNames) {
		queryResult.Error = "updateHostListForParty function: HTTP post request parameter issues (additions lists): facebookID array and name array aren't the same length."
		return queryResult
	}
	queryResult = updateHostListForPartyHelper(r, partyID, hostsToAddFacebookIDs, hostsToAddNames, hostsToRemoveFacebookIDs)
	return queryResult
}

func updateHostListForPartyHelper(r *http.Request, partyID string,
	hostsToAddFacebookIDs []string, hostsToAddNames []string, hostsToRemoveFacebookIDs []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false

	var queryResult1 = QueryResult{}
	queryResult1.Succeeded = true
	var askFriendToHostPartyQueryResult = QueryResult{}
	for i := 0; i < len(hostsToAddFacebookIDs); i++ {
		askFriendToHostPartyQueryResult = askFriendToHostPartyWithYouHelper(partyID, hostsToAddFacebookIDs[i], hostsToAddNames[i])
		if askFriendToHostPartyQueryResult.Succeeded == false {
			queryResult1 = convertTwoQueryResultsToOne(askFriendToHostPartyQueryResult, queryResult1)
		}
	}
	var queryResult2 = QueryResult{}
	queryResult2.Succeeded = true
	var removePartyHostQueryResult = QueryResult{}
	for i := 0; i < len(hostsToRemoveFacebookIDs); i++ {
		removePartyHostQueryResult = removePartyHostHelper(partyID, hostsToRemoveFacebookIDs[i])
		if removePartyHostQueryResult.Succeeded == false {
			queryResult2 = convertTwoQueryResultsToOne(removePartyHostQueryResult, queryResult2)
		}
	}

	queryResult = convertTwoQueryResultsToOne(queryResult1, queryResult2)

	if queryResult1.Succeeded == true {
		hostIsMale, _ := strconv.ParseBool(r.Form.Get("isMale"))
		genderString := "her"
		if hostIsMale == true {
			genderString = "him"
		}
		message := r.Form.Get("name") + " wants you to host a party with " + genderString + "."
		sendPushNotificationsQueryResult := createAndSendNotificationsToThesePeople(hostsToAddFacebookIDs, message, partyID)
		queryResult = convertTwoQueryResultsToOne(queryResult, sendPushNotificationsQueryResult)
	}
	return queryResult
}

func updateHostListForBar(r *http.Request, barID string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false
	r.ParseForm()
	var hostsToAddFacebookIDs []string
	var hostsToAddNames []string
	var hostsToRemoveFacebookIDs []string

	if r.Form.Get("hostsToAddFacebookIDs") != "" {
		hostsToAddFacebookIDs = strings.Split(r.Form.Get("hostsToAddFacebookIDs"), ",")
	}
	if r.Form.Get("hostsToAddNames") != "" {
		hostsToAddNames = strings.Split(r.Form.Get("hostsToAddNames"), ",")
	}
	if r.Form.Get("hostsToRemoveFacebookIDs") != "" {
		hostsToRemoveFacebookIDs = strings.Split(r.Form.Get("hostsToRemoveFacebookIDs"), ",")
	}
	if len(hostsToAddFacebookIDs) != len(hostsToAddNames) {
		queryResult.Error = "updateHostListForBar function: HTTP post request parameter issues (additions lists): facebookID array and name array aren't the same length."
		return queryResult
	}
	queryResult = updateHostListForBarHelper(r, barID, hostsToAddFacebookIDs, hostsToAddNames, hostsToRemoveFacebookIDs)
	return queryResult
}

func updateHostListForBarHelper(r *http.Request, barID string,
	hostsToAddFacebookIDs []string, hostsToAddNames []string, hostsToRemoveFacebookIDs []string) QueryResult {
	var queryResult = QueryResult{}
	queryResult.Succeeded = false

	var queryResult1 = QueryResult{}
	queryResult1.Succeeded = true
	var askFriendToHostBarQueryResult = QueryResult{}
	for i := 0; i < len(hostsToAddFacebookIDs); i++ {
		askFriendToHostBarQueryResult = askFriendToHostBarWithYouHelper(barID, hostsToAddFacebookIDs[i], hostsToAddNames[i])
		if askFriendToHostBarQueryResult.Succeeded == false {
			queryResult1 = convertTwoQueryResultsToOne(askFriendToHostBarQueryResult, queryResult1)
		}
	}
	var queryResult2 = QueryResult{}
	queryResult2.Succeeded = true
	var removeBarHostQueryResult = QueryResult{}
	for i := 0; i < len(hostsToRemoveFacebookIDs); i++ {
		removeBarHostQueryResult = removeBarHostHelper(barID, hostsToRemoveFacebookIDs[i])
		if removeBarHostQueryResult.Succeeded == false {
			queryResult2 = convertTwoQueryResultsToOne(removeBarHostQueryResult, queryResult2)
		}
	}

	queryResult = convertTwoQueryResultsToOne(queryResult1, queryResult2)

	if queryResult1.Succeeded == true {
		hostIsMale, _ := strconv.ParseBool(r.Form.Get("isMale"))
		genderString := "her"
		if hostIsMale == true {
			genderString = "him"
		}
		message := r.Form.Get("nameOfCreator") + " wants you to host a bar with " + genderString + "."
		sendPushNotificationsQueryResult := createAndSendNotificationsToThesePeople(hostsToAddFacebookIDs, message, barID)
		queryResult = convertTwoQueryResultsToOne(queryResult, sendPushNotificationsQueryResult)
	}
	return queryResult
}

func updateInvitationsListAsHostForPartyHelper(partyID string, numberOfInvitesToGive string,
	additionsListFacebookID []string, additionsListIsMale []bool,
	additionsListName []string, removalsListFacebookID []string) QueryResult {
	// this is a random fbid since it doesn't matter what fbid we send in this case
	var myFacebookID = "-1"
	var isHost = true

	var queryResult = QueryResult{}
	queryResult.Succeeded = false

	var queryResult1 = QueryResult{}
	queryResult1.Succeeded = true
	var inviteFriendQueryResult = QueryResult{}
	for i := 0; i < len(additionsListFacebookID); i++ {
		inviteFriendQueryResult = inviteFriendToPartyHelper(partyID, myFacebookID, isHost, numberOfInvitesToGive, additionsListFacebookID[i], additionsListIsMale[i], additionsListName[i])
		if inviteFriendQueryResult.Succeeded == false {
			queryResult1 = convertTwoQueryResultsToOne(inviteFriendQueryResult, queryResult1)
		}
	}
	var queryResult2 = QueryResult{}
	queryResult2.Succeeded = true
	var removeFriendQueryResult = QueryResult{}
	for i := 0; i < len(removalsListFacebookID); i++ {
		removeFriendQueryResult = removeFriendFromPartyHelper(partyID, removalsListFacebookID[i])
		if removeFriendQueryResult.Succeeded == false {
			queryResult2 = convertTwoQueryResultsToOne(removeFriendQueryResult, queryResult2)
		}
	}

	queryResult = convertTwoQueryResultsToOne(queryResult1, queryResult2)
	return queryResult
}

// Find all of the bars for barIDs passed in.
func getBars(w http.ResponseWriter, r *http.Request) {
	var queryResult = QueryResult{}
	if r.URL.Query().Get("barIDs") == "" {
		queryResult.Succeeded = true
		queryResult.DynamodbCalls = nil
	} else {
		barIDs := strings.Split(r.URL.Query().Get("barIDs"), ",")
		queryResult = getBarsHelper(barIDs)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func getClaimKeyHelper(key string) QueryResult {
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
		queryResult.Error = "getClaimKeyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	getItemInput.SetTableName("BarKey")
	var attributeValue = dynamodb.AttributeValue{}
	attributeValue.SetS(key)
	getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"key": &attributeValue})
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getClaimKeyHelper function: GetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil

	data := getItemOutput.Item
	var barKey BarKey
	jsonErr := dynamodbattribute.UnmarshalMap(data, &barKey)
	if jsonErr != nil {
		queryResult.Error = "getClaimKey function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.Succeeded = true
	if barKey.BarID == "" {
		queryResult.Succeeded = false
		queryResult.Error = "The claim key doesn't exist."
	} else {
		queryResult.Error = barKey.BarID
	}
	return queryResult
}

func getBarKeyHelper(key string) QueryResult {
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
		queryResult.Error = "getBarKeyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	getItemInput.SetTableName("BarKey")
	var attributeValue = dynamodb.AttributeValue{}
	attributeValue.SetS(key)
	getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"key": &attributeValue})
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getBarKeyHelper function: GetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil

	data := getItemOutput.Item
	var barKey BarKey
	jsonErr := dynamodbattribute.UnmarshalMap(data, &barKey)
	if jsonErr != nil {
		queryResult.Error = "getBarKey function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.Succeeded = true
	if barKey.Address == "" {
		queryResult.Succeeded = false
		queryResult.Error = "The bar key doesn't exist."
	} else {
		queryResult.Error = barKey.Address
	}
	return queryResult
}

func deleteBarKeyHelper(key string) QueryResult {
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
		queryResult.Error = "deleteBarKey function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var keyAttributeValue = dynamodb.AttributeValue{}
	keyAttributeValue.SetS(key)
	keyMap["key"] = &keyAttributeValue

	var deleteItemInput = dynamodb.DeleteItemInput{}
	deleteItemInput.SetTableName("BarKey")
	deleteItemInput.SetKey(keyMap)

	_, err2 := getter.DynamoDB.DeleteItem(&deleteItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "deleteBarKeyHelper function: DeleteItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func createBarHelper(barKey string, facebookID string, isMale bool, nameOfCreator string, address string, attendeesMapCleanUpHourInZulu string, barID string, details string, latitude string, longitude string, name string, phoneNumber string, schedule map[string]ScheduleForDay, timeZone string, googlePlaceID string) QueryResult {
	var queryResult QueryResult
	if barKey != "AdminUe28GTttHi3L30Jjd3ILLLAdmin" {
		queryResult = getBarKeyHelper(barKey)
		if queryResult.Succeeded == false {
			return queryResult
		}
		queryResult = deleteBarKeyHelper(barKey)
		if queryResult.Succeeded == false {
			return queryResult
		}
	}
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
	var addressAttributeValue = dynamodb.AttributeValue{}
	var attendeesMapCleanUpHourInZuluAttributeValue = dynamodb.AttributeValue{}
	var barIDAttributeValue = dynamodb.AttributeValue{}
	var detailsAttributeValue = dynamodb.AttributeValue{}
	var latitudeAttributeValue = dynamodb.AttributeValue{}
	var longitudeAttributeValue = dynamodb.AttributeValue{}
	var nameAttributeValue = dynamodb.AttributeValue{}
	var phoneNumberAttributeValue = dynamodb.AttributeValue{}
	var timeZoneAttributeValue = dynamodb.AttributeValue{}
	var googlePlaceIDAttributeValue = dynamodb.AttributeValue{}
	addressAttributeValue.SetS(address)
	attendeesMapCleanUpHourInZuluAttributeValue.SetN(attendeesMapCleanUpHourInZulu)
	barIDAttributeValue.SetS(barID)
	detailsAttributeValue.SetS(details)
	latitudeAttributeValue.SetN(latitude)
	longitudeAttributeValue.SetN(longitude)
	nameAttributeValue.SetS(name)
	phoneNumberAttributeValue.SetS(phoneNumber)
	timeZoneAttributeValue.SetN(timeZone)
	googlePlaceIDAttributeValue.SetS(googlePlaceID)
	expressionValues["address"] = &addressAttributeValue
	expressionValues["attendeesMapCleanUpHourInZulu"] = &attendeesMapCleanUpHourInZuluAttributeValue
	expressionValues["barID"] = &barIDAttributeValue
	expressionValues["details"] = &detailsAttributeValue
	expressionValues["latitude"] = &latitudeAttributeValue
	expressionValues["longitude"] = &longitudeAttributeValue
	expressionValues["name"] = &nameAttributeValue
	expressionValues["phoneNumber"] = &phoneNumberAttributeValue
	expressionValues["timeZone"] = &timeZoneAttributeValue
	expressionValues["googlePlaceID"] = &googlePlaceIDAttributeValue

	// set yourself as an attendee to your own bar so that you can rate it
	attendeesMap := make(map[string]*dynamodb.AttributeValue)
	var attendees = dynamodb.AttributeValue{}
	attendeeMap := make(map[string]*dynamodb.AttributeValue)
	var attendee = dynamodb.AttributeValue{}
	var atBarAttribute = dynamodb.AttributeValue{}
	var isMaleAttribute = dynamodb.AttributeValue{}
	var nameOfCreatorAttribute = dynamodb.AttributeValue{}
	var ratingAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	var timeOfLastKnownLocationAttribute = dynamodb.AttributeValue{}
	var timeOfCheckInAttribute = dynamodb.AttributeValue{}
	var saidThereWasACoverAttribute = dynamodb.AttributeValue{}
	var saidThereWasALineAttribute = dynamodb.AttributeValue{}
	atBarAttribute.SetBOOL(false)
	isMaleAttribute.SetBOOL(isMale)
	nameOfCreatorAttribute.SetS(nameOfCreator)
	ratingAttribute.SetS("None")
	statusAttribute.SetS("Maybe")
	timeLastRatedAttribute.SetS("2000-01-01T00:00:00Z")
	timeOfLastKnownLocationAttribute.SetS("2000-01-01T00:00:00Z")
	timeOfCheckInAttribute.SetS("2000-01-01T00:00:00Z")
	saidThereWasACoverAttribute.SetBOOL(false)
	saidThereWasALineAttribute.SetBOOL(false)

	attendeeMap["atBar"] = &atBarAttribute
	attendeeMap["isMale"] = &isMaleAttribute
	attendeeMap["name"] = &nameOfCreatorAttribute
	attendeeMap["rating"] = &ratingAttribute
	attendeeMap["status"] = &statusAttribute
	attendeeMap["timeLastRated"] = &timeLastRatedAttribute
	attendeeMap["timeOfLastKnownLocation"] = &timeOfLastKnownLocationAttribute
	attendeeMap["timeOfCheckIn"] = &timeOfCheckInAttribute
	attendeeMap["saidThereWasACover"] = &saidThereWasACoverAttribute
	attendeeMap["saidThereWasALine"] = &saidThereWasALineAttribute

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
	var hostStatusAttribute = dynamodb.AttributeValue{}
	isMainHostAttribute.SetBOOL(true)
	hostStatusAttribute.SetS("Accepted")
	nameOfCreatorAttribute.SetS(nameOfCreator)
	hostMap["isMainHost"] = &isMainHostAttribute
	hostMap["name"] = &nameOfCreatorAttribute
	hostMap["status"] = &hostStatusAttribute
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
		queryResult.Error += dynamodbCall.Error
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
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func createPartyHelper(facebookID string, isMale bool, name string, address string, details string, drinksProvided bool, endTime string, feeForDrinks bool, invitesForNewInvitees string, latitude string, longitude string, partyID string, startTime string, title string) QueryResult {
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
	var addressAttributeValue = dynamodb.AttributeValue{}
	var detailsAttributeValue = dynamodb.AttributeValue{}
	var drinksProvidedAttributeValue = dynamodb.AttributeValue{}
	var endTimeAttributeValue = dynamodb.AttributeValue{}
	var feeForDrinksAttributeValue = dynamodb.AttributeValue{}
	var invitesForNewInviteesAttributeValue = dynamodb.AttributeValue{}
	var latitudeAttributeValue = dynamodb.AttributeValue{}
	var longitudeAttributeValue = dynamodb.AttributeValue{}
	var partyIDAttributeValue = dynamodb.AttributeValue{}
	var startTimeAttributeValue = dynamodb.AttributeValue{}
	var titleAttributeValue = dynamodb.AttributeValue{}
	invitesForNewInviteesAttributeValue.SetN(invitesForNewInvitees)
	addressAttributeValue.SetS(address)
	detailsAttributeValue.SetS(details)
	drinksProvidedAttributeValue.SetBOOL(drinksProvided)
	endTimeAttributeValue.SetS(endTime)
	feeForDrinksAttributeValue.SetBOOL(feeForDrinks)
	invitesForNewInviteesAttributeValue.SetN(invitesForNewInvitees)
	latitudeAttributeValue.SetN(latitude)
	longitudeAttributeValue.SetN(longitude)
	partyIDAttributeValue.SetS(partyID)
	startTimeAttributeValue.SetS(startTime)
	titleAttributeValue.SetS(title)
	expressionValues["address"] = &addressAttributeValue
	expressionValues["details"] = &detailsAttributeValue
	expressionValues["drinksProvided"] = &drinksProvidedAttributeValue
	expressionValues["endTime"] = &endTimeAttributeValue
	expressionValues["feeForDrinks"] = &feeForDrinksAttributeValue
	expressionValues["invitesForNewInvitees"] = &invitesForNewInviteesAttributeValue
	expressionValues["latitude"] = &latitudeAttributeValue
	expressionValues["longitude"] = &longitudeAttributeValue
	expressionValues["partyID"] = &partyIDAttributeValue
	expressionValues["startTime"] = &startTimeAttributeValue
	expressionValues["title"] = &titleAttributeValue

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
	var atPartyAttribute = dynamodb.AttributeValue{}
	var timeLastRatedAttribute = dynamodb.AttributeValue{}
	var timeOfLastKnownLocationAttribute = dynamodb.AttributeValue{}
	isMaleAttribute.SetBOOL(isMale)
	nameAttribute.SetS(name)
	numberOfInvitationsLeftAttribute.SetN("0")
	ratingAttribute.SetS("None")
	statusAttribute.SetS("Going")
	atPartyAttribute.SetBOOL(false)
	timeLastRatedAttribute.SetS("2000-01-01T00:00:00Z")
	timeOfLastKnownLocationAttribute.SetS("2000-01-01T00:00:00Z")
	inviteeMap["isMale"] = &isMaleAttribute
	inviteeMap["name"] = &nameAttribute
	inviteeMap["numberOfInvitationsLeft"] = &numberOfInvitationsLeftAttribute
	inviteeMap["rating"] = &ratingAttribute
	inviteeMap["atParty"] = &atPartyAttribute
	inviteeMap["status"] = &statusAttribute
	inviteeMap["timeLastRated"] = &timeLastRatedAttribute
	inviteeMap["timeOfLastKnownLocation"] = &timeOfLastKnownLocationAttribute
	invitee.SetM(inviteeMap)
	inviteesMap[facebookID] = &invitee
	invitees.SetM(inviteesMap)
	expressionValues["invitees"] = &invitees

	hostsMap := make(map[string]*dynamodb.AttributeValue)
	var hosts = dynamodb.AttributeValue{}
	hostMap := make(map[string]*dynamodb.AttributeValue)
	var host = dynamodb.AttributeValue{}
	var isMainHostAttribute = dynamodb.AttributeValue{}
	var hostStatusAttribute = dynamodb.AttributeValue{}
	isMainHostAttribute.SetBOOL(true)
	nameAttribute.SetS(name)
	hostStatusAttribute.SetS("Accepted")
	hostMap["isMainHost"] = &isMainHostAttribute
	hostMap["name"] = &nameAttribute
	hostMap["status"] = &hostStatusAttribute
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
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall

	// Now we need to update the person's information to let them
	//     know that they are invited to this party and that they
	//     are hosting this party.
	expressionAttributeNames := make(map[string]*string)
	invitedTo := "invitedTo"
	var partyHostFor = "partyHostFor"
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
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func updatePartyHelper(address string, details string, drinksProvided bool, endTime string, feeForDrinks bool, invitesForNewInvitees string, latitude string, longitude string, partyID string, startTime string, title string) QueryResult {
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
	var addressString = "address"
	var detailsString = "details"
	var drinksProvidedString = "drinksProvided"
	var endTimeString = "endTime"
	var feeForDrinksString = "feeForDrinks"
	var invitesForNewInviteesString = "invitesForNewInvitees"
	var latitudeString = "latitude"
	var longitudeString = "longitude"
	var startTimeString = "startTime"
	var titleString = "title"
	expressionAttributeNames["#address"] = &addressString
	expressionAttributeNames["#details"] = &detailsString
	expressionAttributeNames["#drinksProvided"] = &drinksProvidedString
	expressionAttributeNames["#endTime"] = &endTimeString
	expressionAttributeNames["#feeForDrinks"] = &feeForDrinksString
	expressionAttributeNames["#invitesForNewInvitees"] = &invitesForNewInviteesString
	expressionAttributeNames["#latitude"] = &latitudeString
	expressionAttributeNames["#longitude"] = &longitudeString
	expressionAttributeNames["#startTime"] = &startTimeString
	expressionAttributeNames["#title"] = &titleString
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var addressAttributeValue = dynamodb.AttributeValue{}
	var detailsAttributeValue = dynamodb.AttributeValue{}
	var drinksProvidedAttributeValue = dynamodb.AttributeValue{}
	var endTimeAttributeValue = dynamodb.AttributeValue{}
	var feeForDrinksAttributeValue = dynamodb.AttributeValue{}
	var invitesForNewInviteesAttributeValue = dynamodb.AttributeValue{}
	var latitudeAttributeValue = dynamodb.AttributeValue{}
	var longitudeAttributeValue = dynamodb.AttributeValue{}
	var startTimeAttributeValue = dynamodb.AttributeValue{}
	var titleAttributeValue = dynamodb.AttributeValue{}
	addressAttributeValue.SetS(address)
	detailsAttributeValue.SetS(details)
	drinksProvidedAttributeValue.SetBOOL(drinksProvided)
	endTimeAttributeValue.SetS(endTime)
	feeForDrinksAttributeValue.SetBOOL(feeForDrinks)
	invitesForNewInviteesAttributeValue.SetN(invitesForNewInvitees)
	latitudeAttributeValue.SetN(latitude)
	longitudeAttributeValue.SetN(longitude)
	startTimeAttributeValue.SetS(startTime)
	titleAttributeValue.SetS(title)
	expressionValuePlaceholders[":address"] = &addressAttributeValue
	expressionValuePlaceholders[":details"] = &detailsAttributeValue
	expressionValuePlaceholders[":drinksProvided"] = &drinksProvidedAttributeValue
	expressionValuePlaceholders[":endTime"] = &endTimeAttributeValue
	expressionValuePlaceholders[":feeForDrinks"] = &feeForDrinksAttributeValue
	expressionValuePlaceholders[":invitesForNewInvitees"] = &invitesForNewInviteesAttributeValue
	expressionValuePlaceholders[":latitude"] = &latitudeAttributeValue
	expressionValuePlaceholders[":longitude"] = &longitudeAttributeValue
	expressionValuePlaceholders[":startTime"] = &startTimeAttributeValue
	expressionValuePlaceholders[":title"] = &titleAttributeValue
	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #address=:address, #details=:details, #drinksProvided=:drinksProvided, #endTime=:endTime, #feeForDrinks=:feeForDrinks, #invitesForNewInvitees=:invitesForNewInvitees, #latitude=:latitude, #longitude=:longitude, #startTime=:startTime, #title=:title"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "updatePartyHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func updateBarHelper(address string, attendeesMapCleanUpHourInZulu string, barID string, details string, latitude string, longitude string, name string, phoneNumber string, schedule map[string]ScheduleForDay, timeZone string) QueryResult {
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
	var addressString = "address"
	var attendeesMapCleanUpHourInZuluString = "attendeesMapCleanUpHourInZulu"
	var detailsString = "details"
	var latitudeString = "latitude"
	var longitudeString = "longitude"
	var nameString = "name"
	var phoneNumberString = "phoneNumber"
	var scheduleString = "schedule"
	var timeZoneString = "timeZone"
	expressionAttributeNames["#address"] = &addressString
	expressionAttributeNames["#attendeesMapCleanUpHourInZulu"] = &attendeesMapCleanUpHourInZuluString
	expressionAttributeNames["#details"] = &detailsString
	expressionAttributeNames["#latitude"] = &latitudeString
	expressionAttributeNames["#longitude"] = &longitudeString
	expressionAttributeNames["#name"] = &nameString
	expressionAttributeNames["#phoneNumber"] = &phoneNumberString
	expressionAttributeNames["#schedule"] = &scheduleString
	expressionAttributeNames["#timeZone"] = &timeZoneString
	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var addressAttributeValue = dynamodb.AttributeValue{}
	var attendeesMapCleanUpHourInZuluAttributeValue = dynamodb.AttributeValue{}
	var detailsAttributeValue = dynamodb.AttributeValue{}
	var latitudeAttributeValue = dynamodb.AttributeValue{}
	var longitudeAttributeValue = dynamodb.AttributeValue{}
	var nameAttributeValue = dynamodb.AttributeValue{}
	var phoneNumberAttributeValue = dynamodb.AttributeValue{}
	var timeZoneAttributeValue = dynamodb.AttributeValue{}
	addressAttributeValue.SetS(address)
	attendeesMapCleanUpHourInZuluAttributeValue.SetN(attendeesMapCleanUpHourInZulu)
	detailsAttributeValue.SetS(details)
	latitudeAttributeValue.SetN(latitude)
	longitudeAttributeValue.SetN(longitude)
	nameAttributeValue.SetS(name)
	phoneNumberAttributeValue.SetS(phoneNumber)
	timeZoneAttributeValue.SetN(timeZone)
	expressionValuePlaceholders[":address"] = &addressAttributeValue
	expressionValuePlaceholders[":attendeesMapCleanUpHourInZulu"] = &attendeesMapCleanUpHourInZuluAttributeValue
	expressionValuePlaceholders[":details"] = &detailsAttributeValue
	expressionValuePlaceholders[":latitude"] = &latitudeAttributeValue
	expressionValuePlaceholders[":longitude"] = &longitudeAttributeValue
	expressionValuePlaceholders[":name"] = &nameAttributeValue
	expressionValuePlaceholders[":phoneNumber"] = &phoneNumberAttributeValue
	expressionValuePlaceholders[":timeZone"] = &timeZoneAttributeValue

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
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #address=:address, #attendeesMapCleanUpHourInZulu=:attendeesMapCleanUpHourInZulu, #details=:details, #latitude=:latitude, #longitude=:longitude, #name=:name, #phoneNumber=:phoneNumber, #schedule=:schedule, #timeZone=:timeZone"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "updateBarHelper function: UpdateItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
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
	key.SetS(partyID)
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
		queryResult.Error += dynamodbCall.Error
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
		queryResult.Error += dynamodbCall.Error
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
		queryResult.Error += dynamodbCall.Error
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
	key.SetS(barID)
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
		queryResult.Error += dynamodbCall.Error
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
		queryResult.Error += dynamodbCall.Error
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
	key.SetS(partyID)
	keyMap["partyID"] = &key
	getItemInput.SetKey(keyMap)
	getItemInput.SetTableName("Party")
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getParty function: GetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
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
	key.SetS(barID)
	keyMap["barID"] = &key
	getItemInput.SetKey(keyMap)
	getItemInput.SetTableName("Bar")
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getBar function: GetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
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
	key.SetS(partyID)
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
	updateItemInput.UpdateExpression = &updateExpression
	_, updateItemOutputErr := getter.DynamoDB.UpdateItem(&updateItemInput)

	var dynamodbCall = DynamodbCall{}
	if updateItemOutputErr != nil {
		dynamodbCall.Error = "setNumberOfInvitationsLeftForInviteesHelper function: UpdateItem error (probable cause: this partyID isn't in the database). " + updateItemOutputErr.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func askFriendToHostPartyWithYouHelper(partyID string, friendFacebookID string, name string) QueryResult {
	status := "Waiting"
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
		queryResult.Error = "askFriendToHostPartyWithYouHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	hosts := "hosts"
	expressionAttributeNames["#hosts"] = &hosts
	expressionAttributeNames["#friendFacebookID"] = &friendFacebookID

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var host = dynamodb.AttributeValue{}
	hostMap := make(map[string]*dynamodb.AttributeValue)
	var isMainHostAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	isMainHostAttribute.SetBOOL(false)
	nameAttribute.SetS(name)
	statusAttribute.SetS(status)
	hostMap["isMainHost"] = &isMainHostAttribute
	hostMap["name"] = &nameAttribute
	hostMap["status"] = &statusAttribute
	host.SetM(hostMap)
	expressionValuePlaceholders[":host"] = &host

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #hosts.#friendFacebookID=:host"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "askFriendToHostPartyWithYouHelper function: Problem adding host to Party: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the friend's Person information to let them
	//     know that their friend is asking them to be a host of the Party.
	expressionAttributeNames2 := make(map[string]*string)
	partyHostFor := "partyHostFor"
	expressionAttributeNames2["#partyHostFor"] = &partyHostFor
	expressionAttributeNames2["#partyID"] = &partyID
	expressionValuePlaceholders2 := make(map[string]*dynamodb.AttributeValue)
	var partyIDBoolAttribute = dynamodb.AttributeValue{}
	partyIDBoolAttribute.SetBOOL(false)
	expressionValuePlaceholders2[":bool"] = &partyIDBoolAttribute

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(friendFacebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetExpressionAttributeValues(expressionValuePlaceholders2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	updateExpression2 := "SET #partyHostFor.#partyID=:bool"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, err3 := getter.DynamoDB.UpdateItem(&updateItemInput2)

	var dynamodbCall2 = DynamodbCall{}
	if err3 != nil {
		dynamodbCall2.Error = "askFriendToHostPartyWithYouHelper function: Problem adding partyID to Person's partyHostsFor map. " + err3.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func askFriendToHostBarWithYouHelper(barID string, friendFacebookID string, name string) QueryResult {
	status := "Waiting"
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
		queryResult.Error = "askFriendToHostBarWithYouHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var hosts = "hosts"
	expressionAttributeNames["#hosts"] = &hosts
	expressionAttributeNames["#friendFacebookID"] = &friendFacebookID

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var host = dynamodb.AttributeValue{}
	hostMap := make(map[string]*dynamodb.AttributeValue)
	var isMainHostAttribute = dynamodb.AttributeValue{}
	var nameAttribute = dynamodb.AttributeValue{}
	var statusAttribute = dynamodb.AttributeValue{}
	isMainHostAttribute.SetBOOL(false)
	nameAttribute.SetS(name)
	statusAttribute.SetS(status)
	hostMap["isMainHost"] = &isMainHostAttribute
	hostMap["name"] = &nameAttribute
	hostMap["status"] = &statusAttribute
	host.SetM(hostMap)
	expressionValuePlaceholders[":host"] = &host

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "SET #hosts.#friendFacebookID=:host"
	updateItemInput.UpdateExpression = &updateExpression
	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "askFriendToHostBarWithYouHelper function: Problem adding host to Bar: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the friend's Person information to let them
	//     know that their friend is asking them to be a host of the Bar.
	expressionAttributeNames2 := make(map[string]*string)
	barHostFor := "barHostFor"
	expressionAttributeNames2["#barHostFor"] = &barHostFor
	expressionAttributeNames2["#barID"] = &barID
	expressionValuePlaceholders2 := make(map[string]*dynamodb.AttributeValue)
	var barIDBoolAttribute = dynamodb.AttributeValue{}
	barIDBoolAttribute.SetBOOL(false)
	expressionValuePlaceholders2[":bool"] = &barIDBoolAttribute

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(friendFacebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetExpressionAttributeValues(expressionValuePlaceholders2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	updateExpression2 := "SET #barHostFor.#barID=:bool"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, err3 := getter.DynamoDB.UpdateItem(&updateItemInput2)

	var dynamodbCall2 = DynamodbCall{}
	if err3 != nil {
		dynamodbCall2.Error = "askFriendToHostBarWithYouHelper function: Problem adding barID to Person's barHostsFor map. " + err3.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func removePartyHostHelper(partyID string, facebookID string) QueryResult {
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
		queryResult.Error = "removePartyHostHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	hosts := "hosts"
	expressionAttributeNames["#hosts"] = &hosts
	expressionAttributeNames["#facebookID"] = &facebookID

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "REMOVE #hosts.#facebookID"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "removePartyHostHelper function: Problem removing host from Party: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the friend's Person information to let them
	//     know that they aren't a host of the party anymore.
	expressionAttributeNames2 := make(map[string]*string)
	partyHostFor := "partyHostFor"
	expressionAttributeNames2["#partyHostFor"] = &partyHostFor
	expressionAttributeNames2["#partyID"] = &partyID

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(facebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	updateExpression2 := "REMOVE #partyHostFor.#partyID"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, err3 := getter.DynamoDB.UpdateItem(&updateItemInput2)
	var dynamodbCall2 = DynamodbCall{}
	if err3 != nil {
		dynamodbCall2.Error = "removePartyHostHelper function: Problem removing partyID from Person's partyHostsFor map. " + err3.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func removeBarHostHelper(barID string, facebookID string) QueryResult {
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
		queryResult.Error = "removeBarHostHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var hosts = "hosts"
	expressionAttributeNames["#hosts"] = &hosts
	expressionAttributeNames["#facebookID"] = &facebookID

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	updateExpression := "REMOVE #hosts.#facebookID"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "removeBarHostHelper function: Problem removing host from Bar: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the Person's information to let them
	//     know that they are no longer a host of the Bar.
	expressionAttributeNames2 := make(map[string]*string)
	var barHostFor = "barHostFor"
	expressionAttributeNames2["#barHostFor"] = &barHostFor
	expressionAttributeNames2["#barID"] = &barID

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(facebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	updateExpression2 := "REMOVE #barHostFor.#barID"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, err3 := getter.DynamoDB.UpdateItem(&updateItemInput2)
	var dynamodbCall2 = DynamodbCall{}
	if err3 != nil {
		dynamodbCall2.Error = "removeBarHostHelper function: Problem removing barID from Person's barHostsFor map. " + err3.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

// Accept the invitation to host the Party and then invite this Person to the Party.
func acceptInvitationToHostPartyHelper(partyID string, facebookID string, isMale bool, name string) QueryResult {
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
		queryResult.Error = "acceptInvitationToHostPartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var hosts = "hosts"
	var status = "status"
	expressionAttributeNames["#hosts"] = &hosts
	expressionAttributeNames["#facebookID"] = &facebookID
	expressionAttributeNames["#status"] = &status

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var statusAttribute = dynamodb.AttributeValue{}
	statusAttribute.SetS("Accepted")
	expressionValuePlaceholders[":status"] = &statusAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #hosts.#facebookID.#status=:status"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "acceptInvitationToHostPartyHelper function: Problem changing host status in Party table: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the Person's information to let them
	//     know that they are now hosting this Party.
	expressionAttributeNames2 := make(map[string]*string)
	var partyHostFor = "partyHostFor"
	expressionAttributeNames2["#partyHostFor"] = &partyHostFor
	expressionAttributeNames2["#partyID"] = &partyID
	expressionValuePlaceholders2 := make(map[string]*dynamodb.AttributeValue)
	var partyIDBoolAttribute = dynamodb.AttributeValue{}
	partyIDBoolAttribute.SetBOOL(true)
	expressionValuePlaceholders2[":bool"] = &partyIDBoolAttribute

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(facebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetExpressionAttributeValues(expressionValuePlaceholders2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	updateExpression2 := "SET #partyHostFor.#partyID=:bool"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, err3 := getter.DynamoDB.UpdateItem(&updateItemInput2)

	var dynamodbCall2 = DynamodbCall{}
	if err3 != nil {
		dynamodbCall2.Error = "acceptInvitationToHostPartyHelper function: Problem changing Person's host status for partyHostsFor map. " + err3.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}

	inviteFriendQueryResult := inviteFriendToPartyHelper(partyID, facebookID, true, "0", facebookID, isMale, name)
	if inviteFriendQueryResult.Succeeded == false {
		return inviteFriendQueryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

// Accept the invitation to host the Bar
func acceptInvitationToHostBarHelper(barID string, facebookID string) QueryResult {
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
		queryResult.Error = "acceptInvitationToHostBarHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var hosts = "hosts"
	var status = "status"
	expressionAttributeNames["#hosts"] = &hosts
	expressionAttributeNames["#facebookID"] = &facebookID
	expressionAttributeNames["#status"] = &status

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var statusAttribute = dynamodb.AttributeValue{}
	statusAttribute.SetS("Accepted")
	expressionValuePlaceholders[":status"] = &statusAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	var updateExpression = "SET #hosts.#facebookID.#status=:status"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "acceptInvitationToHostBarHelper function: Problem changing host status in Bar table: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the Person's information to let them
	//     know that they are now hosting this Bar.
	expressionAttributeNames2 := make(map[string]*string)
	var barHostFor = "barHostFor"
	expressionAttributeNames2["#barHostFor"] = &barHostFor
	expressionAttributeNames2["#barID"] = &barID
	expressionValuePlaceholders2 := make(map[string]*dynamodb.AttributeValue)
	var barIDBoolAttribute = dynamodb.AttributeValue{}
	barIDBoolAttribute.SetBOOL(true)
	expressionValuePlaceholders2[":bool"] = &barIDBoolAttribute

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(facebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetExpressionAttributeValues(expressionValuePlaceholders2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	var updateExpression2 = "SET #barHostFor.#barID=:bool"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, err3 := getter.DynamoDB.UpdateItem(&updateItemInput2)

	var dynamodbCall2 = DynamodbCall{}
	if err3 != nil {
		dynamodbCall2.Error = "acceptInvitationToHostBarHelper function: Problem changing Person's host status for barHostsFor map. " + err3.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

// Decline the invitation to host the Party.
func declineInvitationToHostPartyHelper(partyID string, facebookID string) QueryResult {
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
		queryResult.Error = "declineInvitationToHostPartyHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var hosts = "hosts"
	var status = "status"
	expressionAttributeNames["#hosts"] = &hosts
	expressionAttributeNames["#facebookID"] = &facebookID
	expressionAttributeNames["#status"] = &status

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var statusAttribute = dynamodb.AttributeValue{}
	statusAttribute.SetS("Declined")
	expressionValuePlaceholders[":status"] = &statusAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(partyID)
	keyMap["partyID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Party")
	updateExpression := "SET #hosts.#facebookID.#status=:status"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "declineInvitationToHostPartyHelper function: Problem changing host status in Party table: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the Person's information to let them
	//     know that they aren't hosting this Party.
	expressionAttributeNames2 := make(map[string]*string)
	var partyHostFor = "partyHostFor"
	expressionAttributeNames2["#partyHostFor"] = &partyHostFor
	expressionAttributeNames2["#partyID"] = &partyID

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(facebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	updateExpression2 := "REMOVE #partyHostFor.#partyID"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, err3 := getter.DynamoDB.UpdateItem(&updateItemInput2)

	var dynamodbCall2 = DynamodbCall{}
	if err3 != nil {
		dynamodbCall2.Error = "declineInvitationToHostPartyHelper function: Problem removing party in partyHostsFor map in Bar table. " + err3.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

// Decline the invitation to host the Bar
func declineInvitationToHostBarHelper(barID string, facebookID string) QueryResult {
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
		queryResult.Error = "declineInvitationToHostBarHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	expressionAttributeNames := make(map[string]*string)
	var hosts = "hosts"
	var status = "status"
	expressionAttributeNames["#hosts"] = &hosts
	expressionAttributeNames["#facebookID"] = &facebookID
	expressionAttributeNames["#status"] = &status

	expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
	var statusAttribute = dynamodb.AttributeValue{}
	statusAttribute.SetS("Declined")
	expressionValuePlaceholders[":status"] = &statusAttribute

	keyMap := make(map[string]*dynamodb.AttributeValue)
	var key = dynamodb.AttributeValue{}
	key.SetS(barID)
	keyMap["barID"] = &key

	var updateItemInput = dynamodb.UpdateItemInput{}
	updateItemInput.SetExpressionAttributeNames(expressionAttributeNames)
	updateItemInput.SetExpressionAttributeValues(expressionValuePlaceholders)
	updateItemInput.SetKey(keyMap)
	updateItemInput.SetTableName("Bar")
	var updateExpression = "SET #hosts.#facebookID.#status=:status"
	updateItemInput.UpdateExpression = &updateExpression

	_, err2 := getter.DynamoDB.UpdateItem(&updateItemInput)
	var dynamodbCall = DynamodbCall{}
	queryResult.DynamodbCalls[0] = dynamodbCall
	if err2 != nil {
		dynamodbCall.Error = "declineInvitationToHostBarHelper function: Problem changing host status in Bar table: " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}

	// Now we need to update the Person's information to let them
	//     know that they aren't hosting this Bar.
	expressionAttributeNames2 := make(map[string]*string)
	var barHostFor = "barHostFor"
	expressionAttributeNames2["#barHostFor"] = &barHostFor
	expressionAttributeNames2["#barID"] = &barID

	keyMap2 := make(map[string]*dynamodb.AttributeValue)
	var key2 = dynamodb.AttributeValue{}
	key2.SetS(facebookID)
	keyMap2["facebookID"] = &key2

	var updateItemInput2 = dynamodb.UpdateItemInput{}
	updateItemInput2.SetExpressionAttributeNames(expressionAttributeNames2)
	updateItemInput2.SetKey(keyMap2)
	updateItemInput2.SetTableName("Person")
	var updateExpression2 = "REMOVE #barHostFor.#barID"
	updateItemInput2.UpdateExpression = &updateExpression2
	_, err3 := getter.DynamoDB.UpdateItem(&updateItemInput2)

	var dynamodbCall2 = DynamodbCall{}
	if err3 != nil {
		dynamodbCall2.Error = "declineInvitationToHostBarHelper function: Problem removing bar in barHostsFor map in Person table. " + err3.Error()
		dynamodbCall2.Succeeded = false
		queryResult.DynamodbCalls[1] = dynamodbCall2
		queryResult.Error += dynamodbCall2.Error
		return queryResult
	}
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func getBarsHelper(barIDs []string) QueryResult {
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
		queryResult.Error = "getBarsHelper function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var batchGetItemInput = dynamodb.BatchGetItemInput{}
	attributesAndValues := make([]map[string]*dynamodb.AttributeValue, len(barIDs))
	for i := 0; i < len(barIDs); i++ {
		var attributeValue = dynamodb.AttributeValue{}
		attributeValue.SetS(barIDs[i])
		attributesAndValues[i] = make(map[string]*dynamodb.AttributeValue)
		attributesAndValues[i]["barID"] = &attributeValue
	}
	var keysAndAttributes dynamodb.KeysAndAttributes
	keysAndAttributes.SetKeys(attributesAndValues)
	requestedItems := make(map[string]*dynamodb.KeysAndAttributes)
	requestedItems["Bar"] = &keysAndAttributes
	batchGetItemInput.SetRequestItems(requestedItems)
	batchGetItemOutput, err2 := getter.DynamoDB.BatchGetItem(&batchGetItemInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getBarsHelper function: BatchGetItem error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := batchGetItemOutput.Responses
	bars := make([]BarData, len(barIDs))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data["Bar"], &bars)
	if jsonErr != nil {
		queryResult.Error = "getBarsHelper function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.Bars = bars
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}
