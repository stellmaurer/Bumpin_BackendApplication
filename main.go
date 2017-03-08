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
	// curl "http://localhost:8080/createParty?facebookID=000037&isMale=true&name=Stephen%20Ellmaurer&addressLine1=blah1&addressLine2=blah2&city=Fox%20Point&country=United%20States&details=none&drinksProvided=true&endTime=02/02/2017%2002:00:00&feeForDrinks=false&invitesForNewInvitees=4&latitude=-34&longitude=12&startTime=02/01/2017%2022:00:00&stateProvinceRegion=Wisconsin&title=Badger%20Bash&zipCode=53217"
	http.HandleFunc("/createParty", createParty)
	// curl "http://localhost:8080/deleteParty?partyID=11154013587666973726"
	http.HandleFunc("/deleteParty", deleteParty)
	// curl "http://localhost:8080/updateParty?addressLine1=blah1&addressLine2=blah2&city=Fox%20Point&country=United%20States&details=none&drinksProvided=true&endTime=02/02/2017%2002:00:00&feeForDrinks=false&invitesForNewInvitees=4&latitude=-34&longitude=12&partyID=3&startTime=02/01/2017%2022:00:00&stateProvinceRegion=Wisconsin&title=Badger%20Bash&zipCode=53217"
	http.HandleFunc("/updateParty", updateParty)

	/*
		Delete the app queries
	*/
	// curl "http://localhost:8080/deletePerson?facebookID=27"
	http.HandleFunc("/deletePerson", deletePerson)

	/*
		Login queries
	*/
	// curl "http://localhost:8080/createOrUpdatePerson?facebookID=27&isMale=true&name=Zander%20Cage"
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
