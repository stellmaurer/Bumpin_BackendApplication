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
	// curl "http://localhost:8080/myParties?partyID=1"
	http.HandleFunc("/myParties", myParties)
	// curl "http://localhost:8080/barsCloseToMe?latitude=37&longitude=-123"
	http.HandleFunc("/barsCloseToMe", barsCloseToMe)
	// curl "http://localhost:8080/changeAttendanceStatusToParty?partyID=1&facebookID=3303&status=M"
	http.HandleFunc("/changeAttendanceStatusToParty", changeAttendanceStatusToParty)
	// curl "http://localhost:8080/changeAttendanceStatusToBar?barID=1&facebookID=370&atBar=true&isMale=true&name=Steve&rating=N&status=G&timeLastRated=blah"
	http.HandleFunc("/changeAttendanceStatusToBar", changeAttendanceStatusToBar)
	// curl "http://localhost:8080/inviteFriendToParty?partyID=1&myFacebookID=3303&friendFacebookID=1111&isMale=true&name=Ilya"
	http.HandleFunc("/inviteFriendToParty", inviteFriendToParty)

	http.ListenAndServe(":8080", nil)
}
