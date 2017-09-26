package main

import "net/http"

func main() {
	//curl http://localhost:5000/createParty -d "facebookID=30&isMale=true&name=Zander%20Dunn&addressLine1=201%20River%20St&city=River%20Hills&country=United%20States&drinksProvided=true&endTime=2017-05-20T00:00:00Z&feeForDrinks=false&invitesForNewInvitees=1&latitude=-59&longitude=42&startTime=2015-01-18T00:00:00Z&stateProvinceRegion=Wisconsin&title=Baseball%20Party&zipCode=53215"
	//curl http://localhost:5000/askFriendToHostPartyWithYou -d "partyID=11154013587666973726&friendFacebookID=010101&name=Gerrard%20Holler"
	//curl http://localhost:5000/inviteFriendToParty -d "partyID=11154013587666973726&myFacebookID=30&isHost=true&numberOfInvitesToGive=4&friendFacebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"

	//curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=12345699033&isMale=false&name=Sarah%20Carlson&addressLine1=201%20University%20Ave&city=Palo%20Alto&country=United%20States&drinksProvided=true&endTime=2017-05-20T09:00:00Z&feeForDrinks=false&invitesForNewInvitees=4&latitude=-122.161694&longitude=37.446406&startTime=2017-05-19T04:30:00Z&stateProvinceRegion=California&title=PartyToEndItAll&zipCode=94045"
	//curl http://localhost:5000/askFriendToHostPartyWithYou -d "partyID=12258969221770119542&friendFacebookID=010101&name=Gerrard%20Holler"
	//curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/inviteFriendToParty -d "partyID=4957609079704639501&myFacebookID=12345699033&isHost=true&numberOfInvitesToGive=4&friendFacebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"
	/*
		Cleanup of parties that have ended
	*/
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/deletePartiesThatHaveExpired
	go http.HandleFunc("/deletePartiesThatHaveExpired", deletePartiesThatHaveExpired)

	/*
		Cleanup of bar attendee list
	*/
	// curl http://localhost:5000/cleanUpAttendeesMapForBarsThatRecentlyClosed
	// 		Cleanup attendee list of any bar that just closed
	go http.HandleFunc("/cleanUpAttendeesMapForBarsThatRecentlyClosed", cleanUpAttendeesMapForBarsThatRecentlyClosed)

	/*
		Just returns the names of the tables that are in the database - used to check health of Elastic Beanstalk servers
	*/
	// curl "http://localhost:5000/tables"
	go http.HandleFunc("/tables", tables)

	/*
		Find tab queries
	*/
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/myParties?partyIDs=3005619273277206682,3581107088474971827"
	go http.HandleFunc("/myParties", getParties)
	// curl "http://localhost:5000/barsCloseToMe?latitude=43&longitude=-92"
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/barsCloseToMe?latitude=43&longitude=-89"
	go http.HandleFunc("/barsCloseToMe", barsCloseToMe)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/changeAttendanceStatusToParty -d "partyID=10583166241324703384&facebookID=10155613117039816&status=Maybe"
	go http.HandleFunc("/changeAttendanceStatusToParty", changeAttendanceStatusToParty)
	// curl http://localhost:5000/changeAttendanceStatusToBar -d "barID=10206058924726147340&facebookID=90&atBar=false&isMale=true&name=Yasuo%20Yi&rating=None&status=Maybe&timeLastRated=2001-01-01T00:00:00Z&timeOfLastKnownLocation=2001-01-01T00:00:00Z"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/changeAttendanceStatusToBar -d "barID=10206058924726147340&facebookID=90&atBar=false&isMale=true&name=Yasuo%20Yi&rating=None&status=Maybe&timeLastRated=2001-01-01T00:00:00Z&timeOfLastKnownLocation=2001-01-01T00:00:00Z"
	// curl http://localhost:5000/changeAttendanceStatusToBar -d "barID=4454693418154387750&facebookID=9321&atBar=false&isMale=false&name=Emily%20Blunt&rating=None&status=Maybe&timeLastRated=2001-01-01T00:00:00Z&timeOfLastKnownLocation=2001-09-01T00:00:00Z"
	go http.HandleFunc("/changeAttendanceStatusToBar", changeAttendanceStatusToBar)
	// curl http://localhost:5000/changeAtPartyStatus -d "partyID=3005619273277206682&facebookID=10155613117039816&atParty=true&timeOfLastKnownLocation=2017-09-04T00:00:00Z"
	go http.HandleFunc("/changeAtPartyStatus", changeAtPartyStatus)
	// curl http://localhost:5000/changeAtBarStatus -d "barID=1&facebookID=370&atBar=false&isMale=true&name=Steve&rating=None&status=Maybe&timeLastRated=2017-03-04T00:00:00Z&timeOfLastKnownLocation=2017-03-04T00:00:00Z"
	// curl http://localhost:5000/changeAtBarStatus -d "barID=3269697223881195499&facebookID=010101&atBar=true&isMale=true&name=Gerrard%20Holler&rating=None&status=Maybe&timeLastRated=2000-01-01T00:00:00Z&timeOfLastKnownLocation=2017-09-04T00:00:00Z"
	go http.HandleFunc("/changeAtBarStatus", changeAtBarStatus)

	// curl http://localhost:5000/inviteFriendToParty -d "partyID=1&myFacebookID=90&isHost=false&numberOfInvitesToGive=4&friendFacebookID=12345699033&isMale=false&name=Sarah%20Carlson"
	go http.HandleFunc("/inviteFriendToParty", inviteFriendToParty)

	/*
		Rate tab queries
	*/
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/rateParty -d "partyID=10583166241324703384&facebookID=10155613117039816&rating=Heating%20Up&timeLastRated=2017-03-04T00:57:00Z&timeOfLastKnownLocation=2017-03-04T00:00:00Z"
	// curl http://localhost:5000/rateParty -d "partyID=3005619273277206682&facebookID=10155613117039816&rating=Heating%20Up&timeLastRated=2017-03-04T00:57:00Z&timeOfLastKnownLocation=2017-03-04T00:00:00Z"
	go http.HandleFunc("/rateParty", rateParty)
	// curl http://localhost:5000/rateBar -d "barID=4454693418154387750&facebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer&rating=Weak&status=Going&timeLastRated=2017-03-04T01:25:00Z&timeOfLastKnownLocation=2017-03-04T01:00:00Z"
	go http.HandleFunc("/rateBar", rateBar)

	/*
		Host tab queries
	*/
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=000037&isMale=true&name=Stephen%20Ellmaurer&addressLine1=blah1&city=Fox%20Point&country=United%20States&drinksProvided=true&endTime=2017-02-03T02:00:00Z&feeForDrinks=false&invitesForNewInvitees=4&latitude=-34&longitude=12&startTime=2017-02-01T22:00:00Z&stateProvinceRegion=Wisconsin&title=Badger%20Bash&zipCode=53217"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=09876&isMale=false&name=Susan%20Ellmaurer&addressLine1=550%20W%20Main%20St&city=Madison&country=United%20States&drinksProvided=true&endTime=2017-04-08T04:00:00Z&feeForDrinks=false&invitesForNewInvitees=4&latitude=43.068284&longitude=-89.391325&startTime=2017-04-08T22:00:00Z&stateProvinceRegion=Wisconsin&title=Lisa's%20Birthday&zipCode=53703"
	go http.HandleFunc("/createParty", createParty)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createBar -d "barKey=gevEdUSDRlVIKPkQ&facebookID=1222222&isMale=false&nameOfCreator=Eva%20Catarina&addressLine1=1421%20Regent%20St&attendeesMapCleanUpHourInZulu=9&city=Madison&country=United%20States&details=Sports%20Bar&latitude=43.067615&longitude=-89.410205&name=SconnieBar&phoneNumber=608-819-8610&stateProvinceRegion=Wisconsin&timeZone=32&zipCode=53711&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createBar -d "barKey=0q8jQL9QYIL2L3jz&facebookID=09876&isMale=false&nameOfCreator=Susan%20Ellmaurer&addressLine1=100%20N%20Los%20Angeles%20Rd&attendeesMapCleanUpHourInZulu=20&city=Los%20Angeles&country=United%20States&details=A%20bar%20for%20moms.&latitude=18&longitude=-129&name=Women&phoneNumber=608-114-2323&stateProvinceRegion=California&timeZone=6&zipCode=99031&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	// curl http://localhost:5000/createBar -d "barKey=xS32Bk4pBAeRQRFF&facebookID=1222222&isMale=false&nameOfCreator=Eva%20Catarina&addressLine1=1421%20Regent%20St&attendeesMapCleanUpHourInZulu=9&city=Madison&country=United%20States&details=Sports%20Bar&latitude=43.067615&longitude=-89.410205&name=SconnieBar&phoneNumber=608-819-8610&stateProvinceRegion=Wisconsin&timeZone=32&zipCode=53711&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	// curl http://localhost:5000/createBar -d "barKey=UgOMLPkCYUZ1fbfe&facebookID=09876&isMale=false&nameOfCreator=Susan%20Ellmaurer&addressLine1=100%20N%20Los%20Angeles%20Rd&attendeesMapCleanUpHourInZulu=20&city=Los%20Angeles&country=United%20States&details=A%20bar%20for%20moms.&latitude=18&longitude=-129&name=Women&phoneNumber=608-114-2323&stateProvinceRegion=California&timeZone=6&zipCode=99031&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	go http.HandleFunc("/createBar", createBar)
	// curl http://localhost:5000/deleteParty -d "partyID=5233516922553495941"
	go http.HandleFunc("/deleteParty", deleteParty)
	// curl http://localhost:5000/deleteBar -d "barID=2629732187453375056"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/deleteBar -d "barID=218591820409326495"
	go http.HandleFunc("/deleteBar", deleteBar)
	// curl http://localhost:5000/updateParty -d "partyID=12258969221770119542&addressLine1=University%20of%20Milwaukee%20Dorms&city=Milwaukee&country=United%20States&details=none&drinksProvided=true&endTime=2016-02-03T02:00:00Z&feeForDrinks=true&invitesForNewInvitees=3&latitude=33&longitude=77&startTime=2016-02-02T19:02:00Z&stateProvinceRegion=Wisconsin&title=Panther%20Game&zipCode=56677"
	go http.HandleFunc("/updateParty", updateParty)
	// curl http://localhost:5000/updateBar -d "barID=17664650520329034593&addressLine1=84%20Strip%20Terrace&attendeesMapCleanUpHourInZulu=21&city=LA&country=United%20States&details=none&latitude=18&longitude=-129&name=Women&phoneNumber=902-555-3001&stateProvinceRegion=CA&timeZone=7&zipCode=99031&Mon=4PM-4AM,3:30AM&Tue=4PM-4AM,3:30AM&Wed=4PM-4AM,3:30AM&Thu=2PM-4AM,3:30AM&Fri=10AM-4AM,3:30AM&Sat=8AM-4AM,3:30AM&Sun=8AM-2AM,1:45AM"
	go http.HandleFunc("/updateBar", updateBar)
	// curl http://localhost:5000/setNumberOfInvitationsLeftForInvitees -d "partyID=1&invitees=1111,3303,4000&invitationsLeft=2,3,4"
	go http.HandleFunc("/setNumberOfInvitationsLeftForInvitees", setNumberOfInvitationsLeftForInvitees)
	// curl http://localhost:5000/askFriendToHostPartyWithYou -d "partyID=1&friendFacebookID=90&name=Yasuo%20Yi"
	go http.HandleFunc("/askFriendToHostPartyWithYou", askFriendToHostPartyWithYou)
	// curl http://localhost:5000/askFriendToHostBarWithYou -d "barID=1&friendFacebookID=90&name=Yasuo%20Yi"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/askFriendToHostBarWithYou -d "barID=820866964051293233&friendFacebookID=13&name=Eva%20Mendes"
	go http.HandleFunc("/askFriendToHostBarWithYou", askFriendToHostBarWithYou)
	// curl http://localhost:5000/removePartyHost -d "partyID=1&facebookID=90"
	go http.HandleFunc("/removePartyHost", removePartyHost)
	// curl http://localhost:5000/removeBarHost -d "barID=1&facebookID=90"
	go http.HandleFunc("/removeBarHost", removeBarHost)
	// curl http://localhost:5000/acceptInvitationToHostParty -d "partyID=1&facebookID=90&isMale=true&name=Yasuo%20Yi"
	go http.HandleFunc("/acceptInvitationToHostParty", acceptInvitationToHostParty)
	// curl http://localhost:5000/acceptInvitationToHostBar -d "barID=1&facebookID=90"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/acceptInvitationToHostBar -d "barID=820866964051293233&facebookID=13"
	go http.HandleFunc("/acceptInvitationToHostBar", acceptInvitationToHostBar)
	// curl http://localhost:5000/declineInvitationToHostParty -d "partyID=1&facebookID=90"
	go http.HandleFunc("/declineInvitationToHostParty", declineInvitationToHostParty)
	// curl http://localhost:5000/declineInvitationToHostBar -d "barID=1&facebookID=90"
	go http.HandleFunc("/declineInvitationToHostBar", declineInvitationToHostBar)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=1&numberOfInvitesToGive=5&additionsListFacebookID=10155613117039816&additionsListIsMale=true&additionsListName=Steve%20Ellmaurer"
	// curl http://localhost:5000/updateInvitationsListAsHostForParty -d "partyID=12258969221770119542&numberOfInvitesToGive=5&additionsListFacebookID=90,12345699033&additionsListIsMale=true,false&additionsListName=Yasuo%20Yi,Sarah%20Carlson&removalsListFacebookID=1222222,7742229197"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=14312634064321518425&numberOfInvitesToGive=5&additionsListFacebookID=10155613117039816&additionsListIsMale=true&additionsListName=Steve%20Ellmaurer"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=14312634064321518425&removalsListFacebookID=10155613117039816"
	go http.HandleFunc("/updateInvitationsListAsHostForParty", updateInvitationsListAsHostForParty)
	// curl "http://localhost:5000/getPartiesImHosting?partyIDs=1,12258969221770119542"
	go http.HandleFunc("/getPartiesImHosting", getParties)
	// curl "http://localhost:5000/getBarsImHosting?barIDs=1,2"
	go http.HandleFunc("/getBarsImHosting", getBars)

	/*
		Delete the app queries
	*/
	// curl http://localhost:5000/deletePerson -d "facebookID=27"
	go http.HandleFunc("/deletePerson", deletePerson)

	/*
		Login queries
	*/
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createOrUpdatePerson -d "facebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createOrUpdatePerson -d "facebookID=5201&isMale=true&name=Zak%20Shires"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createOrUpdatePerson -d "facebookID=09876&isMale=false&name=Susan%20Ellmaurer"
	// curl http://localhost:5000/createOrUpdatePerson -d "facebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"
	go http.HandleFunc("/createOrUpdatePerson", createOrUpdatePerson)
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/getPerson?facebookID=10155613117039816"
	// curl "http://localhost:5000/getPerson?facebookID=10155613117039816"
	go http.HandleFunc("/getPerson", getPerson)

	/*
		More tab queries
	*/
	// curl http://localhost:5000/updateActivityBlockList -d "myFacebookID=1222222&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	go http.HandleFunc("/updateActivityBlockList", updateActivityBlockList)
	// curl http://localhost:5000/updateIgnoreList -d "myFacebookID=12345699033&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	go http.HandleFunc("/updateIgnoreList", updateIgnoreList)

	/*
		Admin queries (after bar owner identity confirmed, create a key for them and send back an email with their key)
	*/
	// curl http://localhost:5000/createBarKeyForBarOwner
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createBarKeyForBarOwner
	go http.HandleFunc("/createBarKeyForBarOwner", createBarKeyForBarOwner)
	// curl http://localhost:5000/getBarKey -d "key="
	go http.HandleFunc("/getBarKey", getBarKey)
	// curl http://localhost:5000/deleteBarKey -d "key="
	go http.HandleFunc("/deleteBarKey", deleteBarKey)

	//curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/deleteBarKey -d "key=utugUFU8rHdSVTdr"
	//curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/myParties?partyIDs=1,2"
	http.ListenAndServe(":5000", nil)

	/*
		Creates a fraternity party in Rochester, New York
		curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=1222222&isMale=false&name=Eva%20Catarina&addressLine1=Wilson%20Blvd&addressLine2=Alpha%20Delta%20Phi%20Fraternity&city=Rochester&country=United%20States&drinksProvided=true&endTime=2017-05-25T04:00:00Z&feeForDrinks=false&invitesForNewInvitees=5&latitude=43.128793&longitude=-77.632627&startTime=2017-05-25T00:00:00Z&stateProvinceRegion=New%20York&title=Grad%20Party&zipCode=14627"
		curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/inviteFriendToParty -d "partyID=3581107088474971827&myFacebookID=1222222&isHost=true&numberOfInvitesToGive=1&friendFacebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"

		Creates a fraternity party in Madison, WI
		curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=010101&isMale=true&name=Gerrard%20Holler&addressLine1=210%20Langdon%20St&city=Madison&country=United%20States&drinksProvided=true&endTime=2017-06-25T04:00:00Z&feeForDrinks=false&invitesForNewInvitees=5&latitude=43.076698&longitude=-89.393055&startTime=2017-06-25T00:00:00Z&stateProvinceRegion=Wisconsin&title=Summer%20Shenanigans&zipCode=53703"
		curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/inviteFriendToParty -d "partyID=10583166241324703384&myFacebookID=010101&isHost=true&numberOfInvitesToGive=1&friendFacebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"
	*/

	/*


		Demo:

		Susan creates a surprise party for my sister's birthday:
			curl http://localhost:5000/createParty -d "facebookID=09876&isMale=false&name=Susan%20Ellmaurer&addressLine1=550%20W%20Main%20St&city=Madison&country=United%20States&drinksProvided=true&endTime=2017-04-08T04:00:00Z&feeForDrinks=false&invitesForNewInvitees=4&latitude=43.068284&longitude=-89.391325&startTime=2017-04-08T22:00:00Z&stateProvinceRegion=Wisconsin&title=Lisa's%20Birthday&zipCode=53703"
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=09876&isMale=false&name=Susan%20Ellmaurer&addressLine1=550%20W%20Main%20St&city=Madison&country=United%20States&drinksProvided=true&endTime=2017-09-25T04:00Z&feeForDrinks=false&invitesForNewInvitees=5&latitude=43.068284&longitude=-89.391325&startTime=2017-09-25T00:00Z&stateProvinceRegion=Wisconsin&title=Lisa's%20Birthday&zipCode=53703"

		******** Find the party ID and update all these next calls with it:

		Susan want's me to be able to edit the party and invite people so she asks me to host the party with her:
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/askFriendToHostPartyWithYou -d "partyID=3005619273277206682&friendFacebookID=10155613117039816&name=Steve%20Ellmaurer"
		I accept the invitation to host the party with her:
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/acceptInvitationToHostParty -d "partyID=3005619273277206682&facebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"
		******** I automatically get invited when I accept the invitation to host the party

		I invite my friend Sarah Carlson to the party and I give her one invitation so that she can bring a friend:
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/inviteFriendToParty -d "partyID=3005619273277206682&myFacebookID=10155613117039816&isHost=true&numberOfInvitesToGive=1&friendFacebookID=12345699033&isMale=false&name=Sarah%20Carlson"
		Susan removes me as a host because I was inviting too many people:
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/removePartyHost -d "partyID=3005619273277206682&facebookID=10155613117039816"
		Susan uninvites me from the party because I made her mad:
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=3005619273277206682&removalsListFacebookID=10155613117039816"
		After I apologize to her for making her mad and beg her to re-invite me, she decides to invite me again:
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=3005619273277206682&numberOfInvitesToGive=5&additionsListFacebookID=10155613117039816&additionsListIsMale=true&additionsListName=Steve%20Ellmaurer"
		Susan then realizes that she forgot to set my invitations to zero, so she updates that now so that I can't invite anyone:
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/setNumberOfInvitationsLeftForInvitees -d "partyID=3005619273277206682&invitees=10155613117039816&invitationsLeft=0"
		After the expected attendance rises too much, Susan realizes she'll need to move the party to her place, so she updates the party details with the new address and while she's at it, she changes the name to mention it will be a surprise.
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateParty -d "partyID=3005619273277206682&addressLine1=6150%20Century%20Ave&addressLine2=Apt%20109&city=Madison&country=United%20States&details=Meet%20us%20out%20back.&drinksProvided=true&endTime=2017-09-25T04:00:00Z&feeForDrinks=false&invitesForNewInvitees=5&latitude=43.106061&longitude=-89.485127&startTime=2017-09-25T00:00:00Z&stateProvinceRegion=Wisconsin&title=Lisa's%20Surprise%20B-day&zipCode=53703"
		The party comes around and it's AWESOME, so I decide to rate it Bumpin.
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/rateParty -d "partyID=3005619273277206682&facebookID=10155613117039816&rating=Bumpin&timeLastRated=2017-05-25T02:00:00Z&timeOfLastKnownLocation=2017-03-04T00:00:00Z"

		******** The rest of the night is wonderful, but it goes by quick and now the party is finished.
		******** The cleanup of finished parties is actually done automatically everday at noon central time,
							but Susan decides to delete the party manually.
		Delete the party:
			curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/deleteParty -d "partyID=3005619273277206682"


	*/
}
