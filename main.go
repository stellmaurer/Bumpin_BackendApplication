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
	// curl "http://localhost:8080/barsCloseToMe?latitude=43&longitude=-89"
	http.HandleFunc("/barsCloseToMe", barsCloseToMe)
	// curl http://localhost:8080/changeAttendanceStatusToParty -d "partyID=1&facebookID=3303&status=M"
	http.HandleFunc("/changeAttendanceStatusToParty", changeAttendanceStatusToParty)
	// curl http://localhost:8080/changeAttendanceStatusToBar -d "barID=1&facebookID=370&isMale=true&name=Steve&rating=N&status=G&timeLastRated=03/04/2017%2000:00:00"
	http.HandleFunc("/changeAttendanceStatusToBar", changeAttendanceStatusToBar)
	// curl http://localhost:8080/inviteFriendToParty -d "partyID=2&myFacebookID=995830&friendFacebookID=12345699033&isMale=false&name=Sarah%20Carlson"
	http.HandleFunc("/inviteFriendToParty", inviteFriendToParty)

	/*
		Rate tab queries
	*/
	// curl http://localhost:8080/rateParty -d "partyID=1&facebookID=1111&rating=H&timeLastRated=03/04/2017%2000:57:00"
	http.HandleFunc("/rateParty", rateParty)
	// curl http://localhost:8080/rateBar -d "barID=1&facebookID=2323&isMale=true&name=Steve&rating=H&timeLastRated=03/04/2017%2001:25:00"
	http.HandleFunc("/rateBar", rateBar)

	/*
		Host tab queries
	*/
	// curl http://localhost:8080/createParty -d "facebookID=000037&isMale=true&name=Stephen%20Ellmaurer&addressLine1=blah1&city=Fox%20Point&country=United%20States&drinksProvided=true&endTime=02/02/2017%2002:00:00&feeForDrinks=false&invitesForNewInvitees=4&latitude=-34&longitude=12&startTime=02/01/2017%2022:00:00&stateProvinceRegion=Wisconsin&title=Badger%20Bash&zipCode=53217"
	http.HandleFunc("/createParty", createParty)
	// curl http://localhost:8080/createBar -d "facebookID=13&isMale=false&nameOfCreator=Eva%20Mendes&addressLine1=200%20Santa%20Flora%20St&city=San%20Francisco&country=United%20States&details=A%20bar%20for%20actors.&latitude=-22&longitude=11&name=Actor%20Paradise&phoneNumber=902-772-0329&stateProvinceRegion=California&zipCode=96040&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	http.HandleFunc("/createBar", createBar)
	// curl http://localhost:8080/deleteParty -d "partyID=7833179233048568588"
	http.HandleFunc("/deleteParty", deleteParty)
	// curl http://localhost:8080/deleteBar -d "barID=11154013587666973726"
	http.HandleFunc("/deleteBar", deleteBar)
	// curl http://localhost:8080/updateParty -d "partyID=12258969221770119542&addressLine1=University%20of%20Milwaukee%20Dorms&city=Milwaukee&country=United%20States&details=none&drinksProvided=true&endTime=02/03/2016%2002:00:00&feeForDrinks=true&invitesForNewInvitees=3&latitude=33&longitude=77&startTime=02/02/2016%2019:02:00&stateProvinceRegion=Wisconsin&title=Panther%20Game&zipCode=56677"
	http.HandleFunc("/updateParty", updateParty)
	// curl http://localhost:8080/updateBar -d "barID=11154013587666973726&addressLine1=84%20Castro%20Street&city=Mountain%20View&country=United%20States&details=none&latitude=35&longitude=-121.5&name=Philly's&phoneNumber=902-555-3001&stateProvinceRegion=CA&zipCode=94043&Mon=4PM-4AM,3:30AM&Tue=4PM-4AM,3:30AM&Wed=4PM-4AM,3:30AM&Thu=2PM-4AM,3:30AM&Fri=10AM-4AM,3:30AM&Sat=8AM-4AM,3:30AM&Sun=8AM-2AM,1:45AM"
	http.HandleFunc("/updateBar", updateBar)
	// curl http://localhost:8080/setNumberOfInvitationsLeftForInvitees -d "partyID=1&invitees=1111,3303,4000&invitationsLeft=2,3,4"
	http.HandleFunc("/setNumberOfInvitationsLeftForInvitees", setNumberOfInvitationsLeftForInvitees)

	/*
		Delete the app queries
	*/
	// curl http://localhost:8080/deletePerson -d "facebookID=27"
	http.HandleFunc("/deletePerson", deletePerson)

	/*
		Login queries
	*/
	// curl http://localhost:8080/createOrUpdatePerson -d "facebookID=1222222&isMale=false&name=Eva%20Catarina"
	http.HandleFunc("/createOrUpdatePerson", createOrUpdatePerson)
	// curl "http://localhost:8080/getPerson?facebookID=27"
	http.HandleFunc("/getPerson", getPerson)

	/*
		More tab queries
	*/
	// curl http://localhost:8080/updateActivityBlockList -d "myFacebookID=1222222&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	http.HandleFunc("/updateActivityBlockList", updateActivityBlockList)
	// curl http://localhost:8080/updateIgnoreList -d "myFacebookID=12345699033&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	http.HandleFunc("/updateIgnoreList", updateIgnoreList)

	http.ListenAndServe(":8080", nil)
}
