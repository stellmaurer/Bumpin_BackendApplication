package main

import "net/http"

func main() {
	// curl "http://localhost:8080/hello"
	http.HandleFunc("/hello", hello)
	// curl "http://localhost:8080/tables"
	http.HandleFunc("/tables", tables)
	// curl "http://localhost:8080/person"
	http.HandleFunc("/person", person)
	// curl "http://localhost:8080/party"
	http.HandleFunc("/party", party)
	// curl "http://localhost:8080/bar"
	http.HandleFunc("/bar", bar)

	/*
		Find tab queries
	*/
	// curl "http://localhost:8080/myParties?partyID=1,2"
	http.HandleFunc("/myParties", myParties)
	// curl "http://localhost:8080/barsCloseToMe?latitude=37&longitude=-123"
	http.HandleFunc("/barsCloseToMe", barsCloseToMe)
	// curl "http://localhost:8080/changeAttendanceStatusToParty?partyID=1&facebookID=3303&status=M"
	http.HandleFunc("/changeAttendanceStatusToParty", changeAttendanceStatusToParty)
	// curl "http://localhost:8080/changeAttendanceStatusToBar?barID=1&facebookID=370&isMale=true&name=Steve&rating=N&status=G&timeLastRated=blah"
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

	http.ListenAndServe(":8080", nil)
}
