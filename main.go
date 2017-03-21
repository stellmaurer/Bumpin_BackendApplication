package main

import "net/http"

func main() {

	/*
		ticker := time.NewTicker(time.Second * 5)
		go func() {
			for t := range ticker.C {
				fmt.Println("Tick at", t)
			}
		}()
	*/
	//curl http://localhost:8080/createParty -d "facebookID=30&isMale=true&name=Zander%20Dunn&addressLine1=201%20River%20St&city=River%20Hills&country=United%20States&drinksProvided=true&endTime=2015-01-19T00:00:00Z&feeForDrinks=false&invitesForNewInvitees=1&latitude=-59&longitude=42&startTime=2015-01-18T00:00:00Z&stateProvinceRegion=Wisconsin&title=Baseball%20Party&zipCode=53215"
	//curl http://localhost:8080/askFriendToHostPartyWithYou -d "partyID=11154013587666973726&friendFacebookID=010101&name=Gerrard%20Holler"
	//curl http://localhost:8080/inviteFriendToParty -d "partyID=11154013587666973726&myFacebookID=30&isHost=true&numberOfInvitesToGive=4&friendFacebookID=1222222&isMale=false&name=Eva%20Catarina"

	//curl http://localhost:8080/createParty -d "facebookID=12345699033&isMale=false&name=Sarah%20Carlson&addressLine1=201%20University%20Ave&city=Madison&country=United%20States&drinksProvided=true&endTime=2017-03-16T09:00:00Z&feeForDrinks=false&invitesForNewInvitees=4&latitude=-50&longitude=50&startTime=2017-03-16T04:30:00Z&stateProvinceRegion=Wisconsin&title=PartyToEndItAll&zipCode=53705"
	//curl http://localhost:8080/askFriendToHostPartyWithYou -d "partyID=12258969221770119542&friendFacebookID=010101&name=Gerrard%20Holler"
	//curl http://localhost:8080/inviteFriendToParty -d "partyID=12258969221770119542&myFacebookID=30&isHost=true&numberOfInvitesToGive=4&friendFacebookID=1222222&isMale=false&name=Eva%20Catarina"
	/*
		Cleanup of parties that have ended
	*/
	// curl http://localhost:8080/deletePartiesThatHaveExpired
	go http.HandleFunc("/deletePartiesThatHaveExpired", deletePartiesThatHaveExpired)

	/*
		Starting queries - won't be used in app
	*/
	// curl "http://localhost:8080/tables"
	go http.HandleFunc("/tables", tables)

	/*
		Find tab queries
	*/
	// curl "http://localhost:8080/myParties?partyIDs=1,2"
	go http.HandleFunc("/myParties", getParties)
	// curl "http://localhost:8080/barsCloseToMe?latitude=43&longitude=-89"
	go http.HandleFunc("/barsCloseToMe", barsCloseToMe)
	// curl http://localhost:8080/changeAttendanceStatusToParty -d "partyID=1&facebookID=3303&status=M"
	go http.HandleFunc("/changeAttendanceStatusToParty", changeAttendanceStatusToParty)
	// curl http://localhost:8080/changeAttendanceStatusToBar -d "barID=1&facebookID=370&isMale=true&name=Steve&rating=N&status=G&timeLastRated=03/04/2017%2000:00:00"
	go http.HandleFunc("/changeAttendanceStatusToBar", changeAttendanceStatusToBar)
	// curl http://localhost:8080/inviteFriendToParty -d "partyID=1&myFacebookID=90&isHost=false&numberOfInvitesToGive=4&friendFacebookID=12345699033&isMale=false&name=Sarah%20Carlson"
	go http.HandleFunc("/inviteFriendToParty", inviteFriendToParty)

	/*
		Rate tab queries
	*/
	// curl http://localhost:8080/rateParty -d "partyID=1&facebookID=1111&rating=H&timeLastRated=03/04/2017%2000:57:00"
	go http.HandleFunc("/rateParty", rateParty)
	// curl http://localhost:8080/rateBar -d "barID=1&facebookID=2323&isMale=true&name=Steve&rating=H&timeLastRated=03/04/2017%2001:25:00"
	go http.HandleFunc("/rateBar", rateBar)

	/*
		Host tab queries
	*/
	// curl http://localhost:8080/createParty -d "facebookID=000037&isMale=true&name=Stephen%20Ellmaurer&addressLine1=blah1&city=Fox%20Point&country=United%20States&drinksProvided=true&endTime=02/02/2017%2002:00:00&feeForDrinks=false&invitesForNewInvitees=4&latitude=-34&longitude=12&startTime=02/01/2017%2022:00:00&stateProvinceRegion=Wisconsin&title=Badger%20Bash&zipCode=53217"
	go http.HandleFunc("/createParty", createParty)
	// curl http://localhost:8080/createBar -d "barKey=MZdsZsdAmA5gnHoq&facebookID=9321&isMale=false&nameOfCreator=Emily%20Blunt&addressLine1=100%20N%20Los%20Angeles%20Rd&city=Los%20Angeles&country=United%20States&details=A%20bar%20for%20soldiers.&latitude=18&longitude=-129&name=Edge%20of%20Tomorrow&phoneNumber=620-114-2323&stateProvinceRegion=California&zipCode=99031&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	go http.HandleFunc("/createBar", createBar)
	// curl http://localhost:8080/deleteParty -d "partyID=7833179233048568588"
	go http.HandleFunc("/deleteParty", deleteParty)
	// curl http://localhost:8080/deleteBar -d "barID=2629732187453375056"
	go http.HandleFunc("/deleteBar", deleteBar)
	// curl http://localhost:8080/updateParty -d "partyID=12258969221770119542&addressLine1=University%20of%20Milwaukee%20Dorms&city=Milwaukee&country=United%20States&details=none&drinksProvided=true&endTime=02/03/2016%2002:00:00&feeForDrinks=true&invitesForNewInvitees=3&latitude=33&longitude=77&startTime=02/02/2016%2019:02:00&stateProvinceRegion=Wisconsin&title=Panther%20Game&zipCode=56677"
	go http.HandleFunc("/updateParty", updateParty)
	// curl http://localhost:8080/updateBar -d "barID=11154013587666973726&addressLine1=84%20Castro%20Street&city=Mountain%20View&country=United%20States&details=none&latitude=35&longitude=-121.5&name=Philly's&phoneNumber=902-555-3001&stateProvinceRegion=CA&zipCode=94043&Mon=4PM-4AM,3:30AM&Tue=4PM-4AM,3:30AM&Wed=4PM-4AM,3:30AM&Thu=2PM-4AM,3:30AM&Fri=10AM-4AM,3:30AM&Sat=8AM-4AM,3:30AM&Sun=8AM-2AM,1:45AM"
	go http.HandleFunc("/updateBar", updateBar)
	// curl http://localhost:8080/setNumberOfInvitationsLeftForInvitees -d "partyID=1&invitees=1111,3303,4000&invitationsLeft=2,3,4"
	go http.HandleFunc("/setNumberOfInvitationsLeftForInvitees", setNumberOfInvitationsLeftForInvitees)
	// curl http://localhost:8080/askFriendToHostPartyWithYou -d "partyID=1&friendFacebookID=90&name=Yasuo%20Yi"
	go http.HandleFunc("/askFriendToHostPartyWithYou", askFriendToHostPartyWithYou)
	// curl http://localhost:8080/askFriendToHostBarWithYou -d "barID=1&friendFacebookID=90&name=Yasuo%20Yi"
	go http.HandleFunc("/askFriendToHostBarWithYou", askFriendToHostBarWithYou)
	// curl http://localhost:8080/removePartyHost -d "partyID=1&facebookID=90"
	go http.HandleFunc("/removePartyHost", removePartyHost)
	// curl http://localhost:8080/removeBarHost -d "barID=1&facebookID=90"
	go http.HandleFunc("/removeBarHost", removeBarHost)
	// curl http://localhost:8080/acceptInvitationToHostParty -d "partyID=1&facebookID=90&isMale=true&name=Yasuo%20Yi"
	go http.HandleFunc("/acceptInvitationToHostParty", acceptInvitationToHostParty)
	// curl http://localhost:8080/acceptInvitationToHostBar -d "barID=1&facebookID=90"
	go http.HandleFunc("/acceptInvitationToHostBar", acceptInvitationToHostBar)
	// curl http://localhost:8080/declineInvitationToHostParty -d "partyID=1&facebookID=90"
	go http.HandleFunc("/declineInvitationToHostParty", declineInvitationToHostParty)
	// curl http://localhost:8080/declineInvitationToHostBar -d "barID=1&facebookID=90"
	go http.HandleFunc("/declineInvitationToHostBar", declineInvitationToHostBar)
	// curl http://localhost:8080/updateInvitationsListAsHostForParty -d "partyID=12258969221770119542&numberOfInvitesToGive=5&additionsListFacebookID=90,12345699033&additionsListIsMale=true,false&additionsListName=Yasuo%20Yi,Sarah%20Carlson&removalsListFacebookID=1222222,7742229197"
	go http.HandleFunc("/updateInvitationsListAsHostForParty", updateInvitationsListAsHostForParty)
	// curl "http://localhost:8080/getPartiesImHosting?partyIDs=1,12258969221770119542"
	go http.HandleFunc("/getPartiesImHosting", getParties)
	// curl "http://localhost:8080/getBarsImHosting?barIDs=1,2"
	go http.HandleFunc("/getBarsImHosting", getBars)

	/*
		Delete the app queries
	*/
	// curl http://localhost:8080/deletePerson -d "facebookID=27"
	go http.HandleFunc("/deletePerson", deletePerson)

	/*
		Login queries
	*/
	// curl http://localhost:8080/createOrUpdatePerson -d "facebookID=1222222&isMale=false&name=Eva%20Catarina"
	go http.HandleFunc("/createOrUpdatePerson", createOrUpdatePerson)
	// curl "http://localhost:8080/getPerson?facebookID=27"
	go http.HandleFunc("/getPerson", getPerson)

	/*
		More tab queries
	*/
	// curl http://localhost:8080/updateActivityBlockList -d "myFacebookID=1222222&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	go http.HandleFunc("/updateActivityBlockList", updateActivityBlockList)
	// curl http://localhost:8080/updateIgnoreList -d "myFacebookID=12345699033&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	go http.HandleFunc("/updateIgnoreList", updateIgnoreList)

	/*
		Admin queries (after bar owner identity confirmed, create a key for them and send back an email with their key)
	*/
	// curl http://localhost:8080/createBarKeyForBarOwner
	go http.HandleFunc("/createBarKeyForBarOwner", createBarKeyForBarOwner)
	// curl http://localhost:8080/getBarKey -d "key="
	go http.HandleFunc("/getBarKey", getBarKey)
	// curl http://localhost:8080/deleteBarKey -d "key="
	go http.HandleFunc("/deleteBarKey", deleteBarKey)

	http.ListenAndServe(":8080", nil)
}
