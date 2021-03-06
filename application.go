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

import "net/http"

func main() {

	/*
		Daily notifications
	*/
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/sendGoingOutStatusNotificationToPeopleWhoHaveFriendsGoingOutAndHaveALocalTimeEqualToSevenPM
	go http.HandleFunc("/sendGoingOutStatusNotificationToPeopleWhoHaveFriendsGoingOutAndHaveALocalTimeEqualToSevenPM", sendGoingOutStatusNotificationToPeopleWhoHaveFriendsGoingOutAndHaveALocalTimeEqualToSevenPM)

	/*
		Push Notification Testing
	*/
	// curl http://localhost:5000/testiOSPushNotification
	// go http.HandleFunc("/testiOSPushNotification", testiOSPushNotification)
	// curl http://localhost:5000/testAndroidPushNotification
	//go http.HandleFunc("/testAndroidPushNotification", testAndroidPushNotification)
	// curl http://localhost:5000/testSendiOSPushNotification -d "deviceToken=fd21a7d4f1da491ab1b8e817fbd5fe602264ae9a5f9ce27780c8104c132bd891"
	go http.HandleFunc("/testSendiOSPushNotification", testSendiOSPushNotification)
	// curl http://localhost:5000/testSendAndroidPushNotification
	go http.HandleFunc("/testSendAndroidPushNotification", testSendAndroidPushNotification)
	// curl http://localhost:5000/testCreateNotification -d "receiverFacebookID=10154326505409816&message=Hello&partyOrBarID=1"
	// go http.HandleFunc("/testCreateNotification", testCreateNotification)
	// curl http://localhost:5000/testGetPeople
	// go http.HandleFunc("/testGetPeople", testGetPeople)
	// curl http://localhost:5000/testCreateAndSendNotificationsToThesePeople
	// go http.HandleFunc("/testCreateAndSendNotificationsToThesePeople", testCreateAndSendNotificationsToThesePeople)
	// curl http://localhost:5000/getNotificationsForPerson -d "facebookID=10154326505409816"
	go http.HandleFunc("/getNotificationsForPerson", getNotificationsForPerson)
	// curl http://localhost:5000/markNotificationAsSeen -d "notificationID=7816555614368222646"
	go http.HandleFunc("/markNotificationAsSeen", markNotificationAsSeen)
	// curl http://localhost:5000/deleteNotification -d "notificationID=10084172745335654142"
	go http.HandleFunc("/deleteNotification", deleteNotification)

	// curl http://localhost:5000/clearOutstandingNotificationCountForPerson -d "facebookID=10154326505409816"
	go http.HandleFunc("/clearOutstandingNotificationCountForPerson", clearOutstandingNotificationCountForPerson)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/incrementOutstandingNotificationCountForPerson -d "facebookID=10154326505409816"
	go http.HandleFunc("/incrementOutstandingNotificationCountForPerson", incrementOutstandingNotificationCountForPerson)

	//curl http://localhost:5000/createParty -d "facebookID=30&isMale=true&name=Zander%20Dunn&address=201%20River%20St%20River%20Hills%20WI%2053215&drinksProvided=true&endTime=2017-05-20T00:00:00Z&feeForDrinks=false&invitesForNewInvitees=1&latitude=-59&longitude=42&startTime=2015-01-18T00:00:00Z&title=Baseball%20Party"
	//curl http://localhost:5000/askFriendToHostPartyWithYou -d "partyID=11154013587666973726&friendFacebookID=010101&name=Gerrard%20Holler"
	//curl http://localhost:5000/inviteFriendToParty -d "partyID=11154013587666973726&myFacebookID=30&isHost=true&numberOfInvitesToGive=4&friendFacebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"

	//curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=12345699033&isMale=false&name=Sarah%20Carlson&address=201%20University%20Ave%20Palo%20Alto%20CA%2094045&drinksProvided=true&endTime=2017-05-20T09:00:00Z&feeForDrinks=false&invitesForNewInvitees=4&latitude=-122.161694&longitude=37.446406&startTime=2017-05-19T04:30:00Z&title=PartyToEndItAll"
	//curl http://localhost:5000/askFriendToHostPartyWithYou -d "partyID=12258969221770119542&friendFacebookID=010101&name=Gerrard%20Holler"
	//curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/inviteFriendToParty -d "partyID=4957609079704639501&myFacebookID=12345699033&isHost=true&numberOfInvitesToGive=4&friendFacebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer"
	/*
		Cleanup of parties that have ended
	*/
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/deletePartiesThatHaveExpired
	// curl http://localhost:5000/deletePartiesThatHaveExpired
	go http.HandleFunc("/deletePartiesThatHaveExpired", deletePartiesThatHaveExpired)

	/*
		Cleanup of bar attendee list
	*/
	// curl http://localhost:5000/cleanUpAttendeesMapForBarsThatRecentlyClosed
	// 		Cleanup attendee list of any bar that just closed
	go http.HandleFunc("/cleanUpAttendeesMapForBarsThatRecentlyClosed", cleanUpAttendeesMapForBarsThatRecentlyClosed)

	/*
		Cleanup of numberOfFriendsThatMightGoOut in Person table
	*/
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/clearNumberOfFriendsThatMightGoOutForPeopleWhereTheirLocalTimeIsMidnight
	go http.HandleFunc("/clearNumberOfFriendsThatMightGoOutForPeopleWhereTheirLocalTimeIsMidnight", clearNumberOfFriendsThatMightGoOutForPeopleWhereTheirLocalTimeIsMidnight)

	/*
		Just returns the names of the tables that are in the database - used to check health of Elastic Beanstalk servers
	*/
	// curl "http://localhost:5000/tables"
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/tables"
	go http.HandleFunc("/tables", tables)

	/*
		Analytics queries
	*/
	// curl http://localhost:5000/logError -d "ID=Find%20Tab&errorType=server&errorDescription=Blah%20blah"
	go http.HandleFunc("/logError", logError)

	/*
		Bug and Feature Requests
	*/
	// curl http://localhost:5000/createBug -d "facebookID=12&description=wow%20it%20works"
	go http.HandleFunc("/createBug", createBug)
	// curl http://localhost:5000/createFeatureRequest -d "facebookID=12&description=wow%20it%20works"
	go http.HandleFunc("/createFeatureRequest", createFeatureRequest)

	/*
		Find tab queries
	*/
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/myParties?partyIDs=3005619273277206682,3581107088474971827"
	go http.HandleFunc("/myParties", getParties)
	// curl "http://localhost:5000/barsCloseToMe?latitude=43.106045&longitude=-89.484873"
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/barsCloseToMe?latitude=43&longitude=-89"
	go http.HandleFunc("/barsCloseToMe", barsCloseToMe)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/changeAttendanceStatusToParty -d "partyID=10583166241324703384&facebookID=10155613117039816&status=Maybe"
	go http.HandleFunc("/changeAttendanceStatusToParty", changeAttendanceStatusToParty)
	// curl http://localhost:5000/changeAttendanceStatusToBar -d "barID=10206058924726147340&facebookID=90&atBar=false&isMale=true&name=Yasuo%20Yi&rating=None&status=Maybe&timeLastRated=2001-01-01T00:00:00Z&timeOfLastKnownLocation=2001-01-01T00:00:00Z"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/changeAttendanceStatusToBar -d "barID=10206058924726147340&facebookID=90&atBar=false&isMale=true&name=Yasuo%20Yi&rating=None&status=Maybe&timeLastRated=2001-01-01T00:00:00Z&timeOfLastKnownLocation=2001-01-01T00:00:00Z"
	// curl http://localhost:5000/changeAttendanceStatusToBar -d "barID=16773670543315734479&facebookID=10154326505409816&atBar=false&isMale=true&name=Steve%20Ellmaurer&rating=None&status=Maybe&timeLastRated=2001-01-01T00:00:00Z&timeOfLastKnownLocation=2001-09-01T00:00:00Z"
	go http.HandleFunc("/changeAttendanceStatusToBar", changeAttendanceStatusToBar)
	// curl http://localhost:5000/changeAtPartyStatus -d "partyID=3005619273277206682&facebookID=10155613117039816&atParty=true&timeOfLastKnownLocation=2017-09-04T00:00:00Z"
	go http.HandleFunc("/changeAtPartyStatus", changeAtPartyStatus)
	// curl http://localhost:5000/changeAtBarStatus -d "barID=1&facebookID=370&atBar=false&isMale=true&name=Steve&rating=None&status=Maybe&timeLastRated=2017-03-04T00:00:00Z&timeOfLastKnownLocation=2017-03-04T00:00:00Z"
	// curl http://localhost:5000/changeAtBarStatus -d "barID=3269697223881195499&facebookID=010101&atBar=true&isMale=true&name=Gerrard%20Holler&rating=None&status=Maybe&timeLastRated=2000-01-01T00:00:00Z&timeOfLastKnownLocation=2017-09-04T00:00:00Z"
	go http.HandleFunc("/changeAtBarStatus", changeAtBarStatus)

	// curl http://localhost:5000/inviteFriendToParty -d "partyID=1&myFacebookID=90&isHost=false&numberOfInvitesToGive=4&friendFacebookID=12345699033&isMale=false&name=Sarah%20Carlson"
	go http.HandleFunc("/inviteFriendToParty", inviteFriendToParty)

	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/sendInvitationsAsGuestOfParty -d "partyID=2040353648901063840&guestName=Steve%20Ellmaurer&guestFacebookID=10154326505409816&additionsListFacebookID=107808443432866&additionsListIsMale=false&additionsListName=Alex%20Datoriod"
	go http.HandleFunc("/sendInvitationsAsGuestOfParty", sendInvitationsAsGuestOfParty)

	/*
		Rate tab queries
	*/
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/rateParty -d "partyID=10583166241324703384&facebookID=10155613117039816&rating=Heat%27n%20Up&timeLastRated=2017-03-04T00:57:00Z&timeOfLastKnownLocation=2017-03-04T00:00:00Z"
	// curl http://localhost:5000/rateParty -d "partyID=3005619273277206682&facebookID=10155613117039816&rating=Heat%27n%20Up&timeLastRated=2017-03-04T00:57:00Z&timeOfLastKnownLocation=2017-03-04T00:00:00Z"
	go http.HandleFunc("/rateParty", rateParty)
	// curl http://localhost:5000/rateBar -d "barID=4454693418154387750&facebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer&rating=Weak&status=Going&timeLastRated=2017-03-04T01:25:00Z&timeOfLastKnownLocation=2017-03-04T01:00:00Z"
	go http.HandleFunc("/rateBar", rateBar)
	// curl http://localhost:5000/clearRatingForParty -d "partyID=3005619273277206682&facebookID=10155613117039816&timeLastRated=2017-03-04T00:57:00Z&timeOfLastKnownLocation=2017-03-04T00:00:00Z"
	go http.HandleFunc("/clearRatingForParty", clearRatingForParty)
	// curl http://localhost:5000/clearRatingForBar -d "barID=4454693418154387750&facebookID=10155613117039816&isMale=true&name=Steve%20Ellmaurer&status=Going&timeLastRated=2017-03-04T01:25:00Z&timeOfLastKnownLocation=2017-03-04T01:00:00Z"
	go http.HandleFunc("/clearRatingForBar", clearRatingForBar)

	/*
		Host tab queries
	*/
	// curl http://localhost:5000/createParty -d "facebookID=123050841787812&isMale=false&name=Melody%20Panil&address=120%20N%20Breese%20Terrace%20Madison%20WI%2053726&drinksProvided=true&endTime=2017-11-27T08:00:00Z&feeForDrinks=true&invitesForNewInvitees=4&latitude=43.070860&longitude=-89.413948&startTime=2017-11-27T01:00:00Z&title=Breese%20Through%20It&additionsListFacebookID=107798829983852,111354699627054&additionsListIsMale=false,false&additionsListName=Nancy%20Greeneescu,Betty%20Chaison&hostListFacebookIDs=122107341882417,115693492525474&hostListNames=Lisa%20Chengberg,Linda%20Qinstein"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=000037&isMale=true&name=Zander%20Blah&address=415%20E%20Bradley%20Rd%20Fox%20Point%20Wisconsin%2053217&drinksProvided=true&endTime=2017-12-03T02:00:00Z&feeForDrinks=false&invitesForNewInvitees=4&latitude=43.161847&longitude=-87.903058&startTime=2017-12-02T22:00:00Z&title=Badger%20Bash"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createParty -d "facebookID=09876&isMale=false&name=Susan%20Ellmaurer&address=550%20W%20Main%20St%20Madison%2053703&drinksProvided=true&endTime=2017-04-08T04:00:00Z&feeForDrinks=false&invitesForNewInvitees=4&latitude=43.068284&longitude=-89.391325&startTime=2017-04-08T22:00:00Z&title=Lisa's%20Birthday"
	go http.HandleFunc("/createParty", createParty)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createBar -d "barKey=gevEdUSDRlVIKPkQ&facebookID=1222222&isMale=false&nameOfCreator=Eva%20Catarina&address=1421%20Regent%20St%20Madison%20WI%2053711&attendeesMapCleanUpHourInZulu=9&details=Sports%20Bar&latitude=43.067615&longitude=-89.410205&name=SconnieBar&phoneNumber=608-819-8610&timeZone=32&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createBar -d "barKey=0q8jQL9QYIL2L3jz&facebookID=09876&isMale=false&nameOfCreator=Susan%20Ellmaurer&address=100%20N%20Los%20Angeles%20Rd%20Los%20Angeles%20California%2099031&attendeesMapCleanUpHourInZulu=20&details=A%20bar%20for%20moms.&latitude=18&longitude=-129&name=Women&phoneNumber=608-114-2323&timeZone=6&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	// curl http://localhost:5000/createBar -d "barKey=xS32Bk4pBAeRQRFF&facebookID=1222222&isMale=false&nameOfCreator=Eva%20Catarina&address=1421%20Regent%20St%20Madison%20WI%2053711&attendeesMapCleanUpHourInZulu=9&details=Sports%20Bar&latitude=43.067615&longitude=-89.410205&name=SconnieBar&phoneNumber=608-819-8610&timeZone=32&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	// curl http://localhost:5000/createBar -d "barKey=UgOMLPkCYUZ1fbfe&facebookID=09876&isMale=false&nameOfCreator=Susan%20Ellmaurer&address=100%20N%20Los%20Angeles%20Rd%20Los%20Angeles%20CA%2099031&attendeesMapCleanUpHourInZulu=20&details=A%20bar%20for%20moms.&latitude=18&longitude=-129&name=Women&phoneNumber=608-114-2323&timeZone=6&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM"
	// curl http://localhost:5000/createBar -d "barKey=0r5qcj3UQHF2elJz&facebookID=111961819566368&isMale=true&nameOfCreator=Will%20Greenart&address=305%20N%20Midvale%20Blvd%20Apt%20D%20Madison%20WI&attendeesMapCleanUpHourInZulu=20&details=A%20bar%20for%20moms.&latitude=43.070011&longitude=-89.450809&name=Madtown%20Moms&phoneNumber=608-114-2323&timeZone=6&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM&hostListFacebookIDs=122107341882417,115693492525474&hostListNames=Lisa%20Chengberg,Linda%20Qinstein"
	go http.HandleFunc("/createBar", createBar)
	// curl http://localhost:5000/claimBar -d "claimKey=FvIUBYWbH6rNd4bt&facebookID=10154326505409816&isMale=true&nameOfCreator=Steve%20Ellmaurer"
	go http.HandleFunc("/claimBar", claimBar)
	// curl http://localhost:5000/deleteParty -d "partyID=5233516922553495941"
	go http.HandleFunc("/deleteParty", deleteParty)
	// curl http://localhost:5000/deleteBar -d "barID=2629732187453375056"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/deleteBar -d "barID="
	go http.HandleFunc("/deleteBar", deleteBar)
	// hostsToAddFacebookIDs&hostsToAddNames&hostsToRemoveFacebookIDs
	// curl http://localhost:5000/updateParty -d "partyID=13078678500578502570&address=8124%20N%20Seneca%20Rd&details=Steve%20=%20The%20Bomb&drinksProvided=true&endTime=2017-12-25T02:00:00Z&feeForDrinks=false&invitesForNewInvitees=3&latitude=43.1647483&longitude=-87.90766209999998&startTime=2017-12-23T19:02:00Z&title=Steves%20DA%20BOMB%20Party&additionsListFacebookID=107798829983852,111354699627054&additionsListIsMale=false,false&additionsListName=Nancy%20Greeneescu,Betty%20Chaison&hostsToAddFacebookIDs=122107341882417,115693492525474&hostsToAddNames=Lisa%20Chengberg,Linda%20Qinstein"
	// curl http://localhost:5000/updateParty -d "partyID=13078678500578502570&address=8124%20N%20Seneca%20Rd&details=Steve%20=%20The%20Bomb&drinksProvided=true&endTime=2017-12-25T02:00:00Z&feeForDrinks=false&invitesForNewInvitees=3&latitude=43.1647483&longitude=-87.90766209999998&startTime=2017-12-23T19:02:00Z&title=Steves%20DA%20BOMB%20Party&removalsListFacebookID=107798829983852,111354699627054&hostsToRemoveFacebookIDs=122107341882417,115693492525474"
	go http.HandleFunc("/updateParty", updateParty)
	// curl http://localhost:5000/updateBar -d "barID=17664650520329034593&address=84%20Strip%20Terrace%20LA%20CA%2099031&attendeesMapCleanUpHourInZulu=21&details=none&latitude=18&longitude=-129&name=Women&phoneNumber=902-555-3001&timeZone=7&Mon=4PM-4AM,3:30AM&Tue=4PM-4AM,3:30AM&Wed=4PM-4AM,3:30AM&Thu=2PM-4AM,3:30AM&Fri=10AM-4AM,3:30AM&Sat=8AM-4AM,3:30AM&Sun=8AM-2AM,1:45AM&hostsToRemoveFacebookIDs=122107341882417,115693492525474"
	// curl http://localhost:5000/updateBar -d "barID=17664650520329034593&address=84%20Strip%20Terrace%20LA%20CA%2099031&attendeesMapCleanUpHourInZulu=21&details=none&latitude=18&longitude=-129&name=Women&phoneNumber=902-555-3001&timeZone=7&Mon=4PM-4AM,3:30AM&Tue=4PM-4AM,3:30AM&Wed=4PM-4AM,3:30AM&Thu=2PM-4AM,3:30AM&Fri=10AM-4AM,3:30AM&Sat=8AM-4AM,3:30AM&Sun=8AM-2AM,1:45AM&hostsToAddFacebookIDs=122107341882417,115693492525474&hostsToAddNames=Lisa%20Chengberg,Linda%20Qinstein"
	// curl http://localhost:5000/updateBar -d "barID=7209710440755890549&facebookID=111961819566368&isMale=true&nameOfCreator=Will%20Greenart&address=305%20N%20Midvale%20Blvd%20Apt%20D%20Madison%20WI&attendeesMapCleanUpHourInZulu=20&details=A%20bar%20for%20moms.&latitude=43.070011&longitude=-89.450809&name=Madtown%20Moms&phoneNumber=608-114-2323&timeZone=6&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM&hostsToRemoveFacebookIDs=122107341882417,115693492525474"
	// curl http://localhost:5000/updateBar -d "barID=7209710440755890549&facebookID=111961819566368&isMale=true&nameOfCreator=Will%20Greenart&address=305%20N%20Midvale%20Blvd%20Apt%20D%20Madison%20WI&attendeesMapCleanUpHourInZulu=20&details=A%20bar%20for%20moms.&latitude=43.070011&longitude=-89.450809&name=Madtown%20Moms&phoneNumber=608-114-2323&timeZone=6&Mon=4PM-2AM,1:45AM&Tue=4PM-2AM,1:45AM&Wed=4PM-2AM,1:45AM&Thu=2PM-2:30AM,2:00AM&Fri=10AM-3AM,2:30AM&Sat=8AM-3AM,2:30AM&Sun=8AM-1AM,12:45AM&hostsToAddFacebookIDs=122107341882417,115693492525474&hostsToAddNames=Lisa%20Chengberg,Linda%20Qinstein"
	go http.HandleFunc("/updateBar", updateBar)
	// curl http://localhost:5000/setNumberOfInvitationsLeftForInvitees -d "partyID=1&invitees=1111,3303,4000&invitationsLeft=2,3,4"
	go http.HandleFunc("/setNumberOfInvitationsLeftForInvitees", setNumberOfInvitationsLeftForInvitees)
	// curl http://localhost:5000/askFriendsToHostPartyWithYou -d "partyID=1&friendFacebookIDList=90&name=Yasuo%20Yi"
	//go http.HandleFunc("/askFriendsToHostPartyWithYou", askFriendsToHostPartyWithYou)
	// curl http://localhost:5000/removePartyHost -d "partyID=1&facebookID=90"
	go http.HandleFunc("/removePartyHost", removePartyHost)
	// curl http://localhost:5000/removeBarHost -d "barID=1&facebookID=90"
	go http.HandleFunc("/removeBarHost", removeBarHost)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/acceptInvitationToHostParty -d "partyID=10278103879012439008&facebookID=10154326505409816&isMale=true&name=Steve%20Ellmaurer"
	// curl http://localhost:5000/acceptInvitationToHostParty -d "partyID=10278103879012439008&facebookID=10154326505409816&isMale=true&name=Steve%20Ellmaurer"
	go http.HandleFunc("/acceptInvitationToHostParty", acceptInvitationToHostParty)
	// curl http://localhost:5000/acceptInvitationToHostBar -d "barID=1&facebookID=90"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/acceptInvitationToHostBar -d "barID=820866964051293233&facebookID=13"
	go http.HandleFunc("/acceptInvitationToHostBar", acceptInvitationToHostBar)
	// curl http://localhost:5000/declineInvitationToHostParty -d "partyID=1&facebookID=90"
	go http.HandleFunc("/declineInvitationToHostParty", declineInvitationToHostParty)
	// curl http://localhost:5000/declineInvitationToHostBar -d "barID=1&facebookID=90"
	go http.HandleFunc("/declineInvitationToHostBar", declineInvitationToHostBar)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=5349936383849162678&numberOfInvitesToGive=10&additionsListFacebookID=107798829983852,111354699627054&additionsListIsMale=false,false&additionsListName=Nancy%20Greeneescu,Betty%20Chaison"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=5349936383849162678&removalsListFacebookID=107798829983852,111354699627054&numberOfInvitesToGive=9&additionsListFacebookID=113057999456597,184484668766597&additionsListIsMale=false,true&additionsListName=Ruth%20Sidhuson,Mike%20Panditman"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=5349936383849162678&numberOfInvitesToGive=9&removalsListFacebookID=113057999456597,184484668766597"
	// curl http://localhost:5000/updateInvitationsListAsHostForParty -d "partyID=5349936383849162678&numberOfInvitesToGive=10&additionsListFacebookID=107798829983852,111354699627054&additionsListIsMale=false,false&additionsListName=Nancy%20Greeneescu,Betty%20Chaison"
	// curl http://localhost:5000/updateInvitationsListAsHostForParty -d "partyID=5349936383849162678&removalsListFacebookID=107798829983852,111354699627054&numberOfInvitesToGive=9&additionsListFacebookID=113057999456597,184484668766597&additionsListIsMale=false,true&additionsListName=Ruth%20Sidhuson,Mike%20Panditman"
	// curl http://localhost:5000/updateInvitationsListAsHostForParty -d "partyID=5349936383849162678&numberOfInvitesToGive=9&removalsListFacebookID=113057999456597,184484668766597"
	// curl http://localhost:5000/updateInvitationsListAsHostForParty -d "partyID=12258969221770119542&numberOfInvitesToGive=5&additionsListFacebookID=90,12345699033&additionsListIsMale=true,false&additionsListName=Yasuo%20Yi,Sarah%20Carlson&removalsListFacebookID=1222222,7742229197"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=14312634064321518425&numberOfInvitesToGive=5&additionsListFacebookID=10155613117039816&additionsListIsMale=true&additionsListName=Steve%20Ellmaurer"
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/updateInvitationsListAsHostForParty -d "partyID=5349936383849162678&removalsListFacebookID=10155613117039816"
	//go http.HandleFunc("/updateInvitationsListAsHostForParty", updateInvitationsListAsHostForParty)
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
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createOrUpdatePerson -d "facebookID=10154326505409816&isMale=true&name=Steve%20Ellmaurer"
	// curl http://localhost:5000/createOrUpdatePerson -d "facebookID=1&isMale=true&name=Zitoli%20Estov&platform=iOS&deviceToken=Unknown"
	go http.HandleFunc("/createOrUpdatePerson", createOrUpdatePerson)
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/getPerson?facebookID=10155613117039816"
	// curl "http://localhost:5000/getPerson?facebookID=10155613117039816"
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/getPerson?facebookID=10154326505409816"
	go http.HandleFunc("/getPerson", getPerson)
	// curl http://localhost:5000/updateWhatGotPersonToDownload -d "facebookID=10154326505409816&whatGotThemToDownload=I%20dont%20know"
	go http.HandleFunc("/updateWhatGotPersonToDownload", updateWhatGotPersonToDownload)

	/*
		More tab queries
	*/
	// curl http://localhost:5000/updateActivityBlockList -d "myFacebookID=1222222&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	go http.HandleFunc("/updateActivityBlockList", updateActivityBlockList)
	// curl http://localhost:5000/updateIgnoreList -d "myFacebookID=12345699033&removalsList=12,010101,90&additionsList=13,12345699033,7742229197"
	go http.HandleFunc("/updateIgnoreList", updateIgnoreList)
	// curl http://localhost:5000/updatePersonStatus -d "facebookID=10154326505409816&goingOut=Yes&timeGoingOutStatusWasSet=2000-01-01T00:00:00Z&manuallySet=Yes"
	go http.HandleFunc("/updatePersonStatus", updatePersonStatus)
	// curl http://localhost:5000/incrementNumberOfFriendsThatMightGoOutForTheseFriends -d "facebookID=10154326505409816&friendFacebookIDs=10155227369101712,1617903301590247,10203989030603248"
	go http.HandleFunc("/incrementNumberOfFriendsThatMightGoOutForTheseFriends", incrementNumberOfFriendsThatMightGoOutForTheseFriends)
	// curl http://localhost:5000/decrementNumberOfFriendsThatMightGoOutForTheseFriends -d "facebookID=10154326505409816&friendFacebookIDs=10155227369101712,1617903301590247,10203989030603248"
	go http.HandleFunc("/decrementNumberOfFriendsThatMightGoOutForTheseFriends", decrementNumberOfFriendsThatMightGoOutForTheseFriends)
	// curl "http://bumpin-env.us-west-2.elasticbeanstalk.com:80/getFriends?facebookIDs=10155227369101712,10203989030603248,102078678319995"
	// curl "http://localhost:5000/getFriends?facebookIDs=10213227731221423,10155513390409021,10108057048841417"
	go http.HandleFunc("/getFriends", getFriends)

	/*
		Admin queries (after bar owner identity confirmed, create a key for them and send back an email with their key)
	*/
	// curl http://localhost:5000/createBarKeyForAddress
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createBarKeyForAddress -d "address=305%20N%20Midvale%20Blvd%20Apt%20D%20Madison%20WI"
	go http.HandleFunc("/createBarKeyForAddress", createBarKeyForAddress)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/createClaimKeyForBar -d "barID=991135472298699438"
	go http.HandleFunc("/createClaimKeyForBar", createClaimKeyForBar)
	// curl http://localhost:5000/getBarKey -d "key="
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/getBarKey -d "key=0yweLNBFuooCgvQD"
	go http.HandleFunc("/getBarKey", getBarKey)
	// curl http://bumpin-env.us-west-2.elasticbeanstalk.com:80/getClaimKey -d "key=0yweLNBFuooCgvQD"
	go http.HandleFunc("/getClaimKey", getClaimKey)
	// curl http://localhost:5000/deleteBarKey -d "key="
	go http.HandleFunc("/deleteBarKey", deleteBarKey)
	// https://maps.googleapis.com/maps/api/place/nearbysearch/json?key=AIzaSyBDGTJegyakdJ3ObWRQfecI9zH_0MyzRhM&location=30.216979,-97.750996&radius=800&type=bar
	// curl http://localhost:5000/populateBarsFromGooglePlacesAPI -d "timeZone=31&attendeesMapCleanUpHourInZulu=10&latitude=40.643200&longitude=-73.949834&radius=1000&squareMiles=8"
	// curl http://localhost:5000/populateBarsFromGooglePlacesAPI -d "timeZone=33&attendeesMapCleanUpHourInZulu=12&latitude=45.014322&longitude=-93.270163&radius=200&squareMiles=60" Minneapolis, MN 37.475278
	// curl http://localhost:5000/populateBarsFromGooglePlacesAPI -d "timeZone=35&attendeesMapCleanUpHourInZulu=14&latitude=37.475278&longitude=-122.151102&radius=400&squareMiles=2" Facebook HQ
	go http.HandleFunc("/populateBarsFromGooglePlacesAPI", populateBarsFromGooglePlacesAPI)
	// curl "http://localhost:5000/deleteBarDuplicates"
	go http.HandleFunc("/deleteBarDuplicates", deleteBarDuplicates)
	// curl "http://localhost:5000/deleteAllBarsThatWereAutoPopulated"
	go http.HandleFunc("/deleteAllBarsThatWereAutoPopulated", deleteAllBarsThatWereAutoPopulated)
	// curl "http://localhost:5000/deleteBarsThatArentOpenAfter1159PMOnFriday"
	go http.HandleFunc("/deleteBarsThatArentOpenAfter1159PMOnFriday", deleteBarsThatArentOpenAfter1159PMOnFriday)
	// curl "http://localhost:5000/getBarsThatBarCreatorNeedsAddedToTheirBarHostForMap"
	go http.HandleFunc("/getBarsThatBarCreatorNeedsAddedToTheirBarHostForMap", getBarsThatBarCreatorNeedsAddedToTheirBarHostForMap)

	http.ListenAndServe(":5000", nil)
}
