package main

import "net/http"

func main() {
	/*
		Starting queries - won't be used in app
	*/
	// curl "http://localhost:8080/tables"
	http.HandleFunc("/tables", tables)

	/*
		Find tab queries
	*/
	// curl "http://localhost:8080/myParties?partyID=1,2"
	http.HandleFunc("/myParties", myParties)
	// curl "http://localhost:8080/barsCloseToMe?latitude=37&longitude=-123"
	http.HandleFunc("/barsCloseToMe", barsCloseToMe)
	// curl "http://localhost:8080/changeAttendanceStatusToParty?partyID=1&facebookID=3303&status=M"
	http.HandleFunc("/changeAttendanceStatusToParty", changeAttendanceStatusToParty)
	// curl "http://localhost:8080/changeAttendanceStatusToBar?barID=1&facebookID=370&isMale=true&name=Steve&rating=N&status=G&timeLastRated=03/04/2017%2000:00:00"
	http.HandleFunc("/changeAttendanceStatusToBar", changeAttendanceStatusToBar)
	// curl "http://localhost:8080/inviteFriendToParty?partyID=2&myFacebookID=995830&friendFacebookID=12345699033&isMale=false&name=Sarah%20Carlson"
	http.HandleFunc("/inviteFriendToParty", inviteFriendToParty)

	/*
		Rate tab queries
	*/
	// curl "http://localhost:8080/rateParty?partyID=1&facebookID=1111&rating=H&timeLastRated=03/04/2017%2000:57:00"
	http.HandleFunc("/rateParty", rateParty)
	// curl "http://localhost:8080/rateBar?barID=1&facebookID=2323&isMale=true&name=Steve&rating=H&timeLastRated=03/04/2017%2001:25:00"
	http.HandleFunc("/rateBar", rateBar)

	/*
		Host tab queries
	*/
	// curl "http://localhost:8080/createParty?facebookID=000037&isMale=true&name=Stephen%20Ellmaurer&addressLine1=blah1&city=Fox%20Point&country=United%20States&drinksProvided=true&endTime=02/02/2017%2002:00:00&feeForDrinks=false&invitesForNewInvitees=4&latitude=-34&longitude=12&startTime=02/01/2017%2022:00:00&stateProvinceRegion=Wisconsin&title=Badger%20Bash&zipCode=53217"
	http.HandleFunc("/createParty", createParty)
	// curl "http://localhost:8080/createBar?facebookID=12345699033&isMale=false&nameOfCreator=Sarah%20Carlson&addressLine1=3001%20University%20St&city=Madison&closingTime=3:00%20AM&country=United%20States&details=A%20bar%20for%20gymnasts&lastCall=2:45%20AM&latitude=-22&longitude=11&name=Gymbar&openAt=4:00%20PM&stateProvinceRegion=Wisconsin&zipCode=53217"
	http.HandleFunc("/createBar", createBar)
	// curl "http://localhost:8080/deleteParty?partyID=7833179233048568588"
	http.HandleFunc("/deleteParty", deleteParty)
	// curl "http://localhost:8080/updateParty?partyID=12258969221770119542&addressLine1=University%20of%20Milwaukee%20Dorms&city=Milwaukee&country=United%20States&details=none&drinksProvided=true&endTime=02/03/2016%2002:00:00&feeForDrinks=true&invitesForNewInvitees=3&latitude=33&longitude=77&startTime=02/02/2016%2019:02:00&stateProvinceRegion=Wisconsin&title=Panther%20Game&zipCode=56677"
	http.HandleFunc("/updateParty", updateParty)

	/*
		Delete the app queries
	*/
	// curl "http://localhost:8080/deletePerson?facebookID=27"
	http.HandleFunc("/deletePerson", deletePerson)

	/*
		Login queries
	*/
	// curl "http://localhost:8080/createOrUpdatePerson?facebookID=000037&isMale=true&name=Zander%20Blah"
	http.HandleFunc("/createOrUpdatePerson", createOrUpdatePerson)
	// curl "http://localhost:8080/getPerson?facebookID=27"
	http.HandleFunc("/getPerson", getPerson)

	/*
		More tab queries
	*/
	// curl "http://localhost:8080/updateBlockList?myFacebookID=1222222&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	http.HandleFunc("/updateBlockList", updateBlockList)

	http.ListenAndServe(":8080", nil)
}
