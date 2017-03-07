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
