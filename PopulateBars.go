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
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func deleteBarDuplicates(w http.ResponseWriter, r *http.Request) {
	var numberOfDuplicates = 0
	barsAlreadyInOurDB := getAllBars().Bars
	barMap := make(map[string]bool)
	for i := 0; i < len(barsAlreadyInOurDB); i++ {
		barExistsAlready := barMap[barsAlreadyInOurDB[i].GooglePlaceID]
		if (barExistsAlready == true) && (barsAlreadyInOurDB[i].GooglePlaceID != "-1") {
			deleteBarHelper(barsAlreadyInOurDB[i].BarID)
			numberOfDuplicates++
		} else {
			barMap[barsAlreadyInOurDB[i].GooglePlaceID] = true
		}
	}
	var message = "Number of bar duplicates deleted = " + strconv.Itoa(numberOfDuplicates)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(message)
}

func populateBarsFromGooglePlacesAPI(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	latitudeOfQuery := r.Form.Get("latitude")
	longitudeOfQuery := r.Form.Get("longitude")
	radius := r.Form.Get("radius")
	squareMiles := r.Form.Get("squareMiles")
	timeZone := r.Form.Get("timeZone")
	attendeesMapCleanUpHourInZulu := r.Form.Get("attendeesMapCleanUpHourInZulu")

	barKey := "AdminUe28GTttHi3L30Jjd3ILLLAdmin"
	facebookID := "184484668766597"
	isMale := true
	nameOfCreator := "Bar Creator"
	details := "None"
	lastCall := "1:45 AM"

	barsDetailed := getBarsFromGooglePlacesAPIWithinThisSquareMileage(latitudeOfQuery, longitudeOfQuery, radius, squareMiles)

	barsAlreadyInOurDB := getBarsWithinThisSquareMileageFromTheseCoordinates(latitudeOfQuery, longitudeOfQuery, squareMiles).Bars
	mapOfBarsAlreadyInOurDB := make(map[string]BarData)
	for i := 0; i < len(barsAlreadyInOurDB); i++ {
		mapOfBarsAlreadyInOurDB[barsAlreadyInOurDB[i].GooglePlaceID] = barsAlreadyInOurDB[i]
	}

	var queryResult QueryResult
	queryResult.Succeeded = true

	for i := 0; i < len(barsDetailed); i++ {
		latitude := strconv.FormatFloat(barsDetailed[i].Geometry.Location.Lat, 'f', -1, 64)
		longitude := strconv.FormatFloat(barsDetailed[i].Geometry.Location.Lng, 'f', -1, 64)
		if latitude == "0" {
			continue
		}
		if longitude == "0" {
			continue
		}
		address := barsDetailed[i].Address
		if address == "" {
			continue
		}
		name := barsDetailed[i].Name
		if name == "" {
			continue
		}
		phoneNumber := barsDetailed[i].PhoneNumber
		if phoneNumber == "" {
			phoneNumber = "N/A"
		}

		schedule := make(map[string]ScheduleForDay)
		if len(barsDetailed[i].OpeningHours.Schedule) >= 7 {
			mondayHours := barsDetailed[i].OpeningHours.Schedule[0]
			if mondayHours == "" {
				mondayHours = "N/A"
			}
			tuesdayHours := barsDetailed[i].OpeningHours.Schedule[1]
			if tuesdayHours == "" {
				tuesdayHours = "N/A"
			}
			wednesdayHours := barsDetailed[i].OpeningHours.Schedule[2]
			if wednesdayHours == "" {
				wednesdayHours = "N/A"
			}
			thursdayHours := barsDetailed[i].OpeningHours.Schedule[3]
			if thursdayHours == "" {
				thursdayHours = "N/A"
			}
			fridayHours := barsDetailed[i].OpeningHours.Schedule[4]
			if fridayHours == "" {
				fridayHours = "N/A"
			}
			saturdayHours := barsDetailed[i].OpeningHours.Schedule[5]
			if saturdayHours == "" {
				saturdayHours = "N/A"
			}
			sundayHours := barsDetailed[i].OpeningHours.Schedule[6]
			if sundayHours == "" {
				sundayHours = "N/A"
			}
			schedule["Monday"] = ScheduleForDay{Open: mondayHours, LastCall: lastCall}
			schedule["Tuesday"] = ScheduleForDay{Open: tuesdayHours, LastCall: lastCall}
			schedule["Wednesday"] = ScheduleForDay{Open: wednesdayHours, LastCall: lastCall}
			schedule["Thursday"] = ScheduleForDay{Open: thursdayHours, LastCall: lastCall}
			schedule["Friday"] = ScheduleForDay{Open: fridayHours, LastCall: lastCall}
			schedule["Saturday"] = ScheduleForDay{Open: saturdayHours, LastCall: lastCall}
			schedule["Sunday"] = ScheduleForDay{Open: sundayHours, LastCall: lastCall}
		} else {
			schedule["Monday"] = ScheduleForDay{Open: "N/A", LastCall: lastCall}
			schedule["Tuesday"] = ScheduleForDay{Open: "N/A", LastCall: lastCall}
			schedule["Wednesday"] = ScheduleForDay{Open: "N/A", LastCall: lastCall}
			schedule["Thursday"] = ScheduleForDay{Open: "N/A", LastCall: lastCall}
			schedule["Friday"] = ScheduleForDay{Open: "N/A", LastCall: lastCall}
			schedule["Saturday"] = ScheduleForDay{Open: "N/A", LastCall: lastCall}
			schedule["Sunday"] = ScheduleForDay{Open: "N/A", LastCall: lastCall}
		}

		existingBar, barExistsAlready := mapOfBarsAlreadyInOurDB[barsDetailed[i].PlaceID]
		if barExistsAlready == false {
			barID := strconv.FormatUint(getRandomID(), 10)
			var createBarQueryResult = createBarHelper(barKey, facebookID, isMale, nameOfCreator, address, attendeesMapCleanUpHourInZulu, barID, details, latitude, longitude, name, phoneNumber, schedule, timeZone, barsDetailed[i].PlaceID)
			if createBarQueryResult.Succeeded == false {
				fmt.Println(barKey, facebookID, isMale, nameOfCreator, address, attendeesMapCleanUpHourInZulu, barID, details, latitude, longitude, name, phoneNumber, schedule, timeZone, barsDetailed[i].PlaceID)
			}
			queryResult = convertTwoQueryResultsToOne(queryResult, createBarQueryResult)
			mapOfBarsAlreadyInOurDB[barsDetailed[i].PlaceID] = getBar(barID).Bars[0]
		} else {
			_, barCreatorIsOwner := existingBar.Hosts["184484668766597"]
			if barCreatorIsOwner == true {
				var updateBarQueryResult = updateBarHelper(address, attendeesMapCleanUpHourInZulu, existingBar.BarID, details, latitude, longitude, name, phoneNumber, schedule, timeZone)
				queryResult = convertTwoQueryResultsToOne(queryResult, updateBarQueryResult)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(queryResult)
}

func getBarsFromGooglePlacesAPIWithinThisSquareMileage(latitudeString string, longitudeString string, radiusString string, squareMileage string) []GooglePlaceDetailed {
	n := calculateWhatNShouldBe(squareMileage, radiusString)
	mapOfBarsAlreadyAdded := make(map[string]bool) // making sure we only return unique bars
	fmt.Println("n = ", n)
	latitude, latitudeErr := strconv.ParseFloat(latitudeString, 64)
	longitude, longitudeErr := strconv.ParseFloat(longitudeString, 64)
	radius, radiusErr := strconv.ParseFloat(radiusString, 64)

	if latitudeErr != nil || longitudeErr != nil || radiusErr != nil {
		return nil
	}

	degreesOfLatitudeWhichEqual1Meter := 0.000009009904681
	degreesOfLongitudeWhichEqual1Meter := (1 / (69.1703 * math.Cos(latitude*0.0174533))) / 1609.34
	startingLatitude := latitude + ((radius * float64(n/2)) * degreesOfLatitudeWhichEqual1Meter)
	startingLongitude := longitude - ((radius * float64(n/2)) * degreesOfLongitudeWhichEqual1Meter)

	var bars []GooglePlaceDetailed
	currentLatitude := startingLatitude
	currentLongitude := startingLongitude
	for i := 0; i < n; i++ {
		degreesOfLongitudeWhichEqual1Meter = (1 / (69.1703 * math.Cos(currentLatitude*0.0174533))) / 1609.34
		for j := 0; j < n; j++ {
			var barsToAdd []GooglePlaceDetailed
			if ((i % 2) == 0) && ((j % 2) == 0) {
				barsToAdd = getBarsFromGooglePlacesAPIWithinThisRadiusAndNotInTheMap(currentLatitude, currentLongitude, radiusString, mapOfBarsAlreadyAdded)
			}
			if ((i % 2) != 0) && ((j % 2) != 0) {
				barsToAdd = getBarsFromGooglePlacesAPIWithinThisRadiusAndNotInTheMap(currentLatitude, currentLongitude, radiusString, mapOfBarsAlreadyAdded)
			}
			bars = append(bars, barsToAdd...)
			currentLongitude = currentLongitude + (radius * degreesOfLongitudeWhichEqual1Meter)
		}
		currentLongitude = startingLongitude
		currentLatitude = currentLatitude - (radius * degreesOfLatitudeWhichEqual1Meter)
	}

	return bars
}

func getBarsFromGooglePlacesAPIWithinThisRadiusAndNotInTheMap(latitudeFloat float64, longitudeFloat float64, radius string, mapOfBarsAlreadyAdded map[string]bool) []GooglePlaceDetailed {
	latitude := strconv.FormatFloat(latitudeFloat, 'f', -1, 64)
	longitude := strconv.FormatFloat(longitudeFloat, 'f', -1, 64)
	var bars []GooglePlace = getPlaceIDsOfBarsFromGooglePlacesAPI(latitude, longitude, radius)

	var barsDetailed []GooglePlaceDetailed

	for i := 0; i < len(bars); i++ {
		var barToAdd = getPlaceDetailsForPlaceID(bars[i].PlaceID)
		barExistsAlready := mapOfBarsAlreadyAdded[barToAdd.PlaceID]
		if barExistsAlready == false {
			mapOfBarsAlreadyAdded[barToAdd.PlaceID] = true
			barsDetailed = append(barsDetailed, barToAdd)
		}
	}
	return barsDetailed
}

func getBarsFromGooglePlacesAPIWithinThisRadius(latitude string, longitude string, radius string) []GooglePlaceDetailed {
	var bars []GooglePlace = getPlaceIDsOfBarsFromGooglePlacesAPI(latitude, longitude, radius)
	var barsDetailed []GooglePlaceDetailed

	for i := 0; i < len(bars); i++ {
		barsDetailed = append(barsDetailed, getPlaceDetailsForPlaceID(bars[i].PlaceID))
	}

	return barsDetailed
}

// https://maps.googleapis.com/maps/api/place/details/json?key=AIzaSyBDGTJegyakdJ3ObWRQfecI9zH_0MyzRhM&placeid=
func getPlaceDetailsForPlaceID(placeID string) GooglePlaceDetailed {
	var result Result

	var query string = "https://maps.googleapis.com/maps/api/place/details/json?key=AIzaSyBDGTJegyakdJ3ObWRQfecI9zH_0MyzRhM&placeid=" + placeID
	resp, err := http.Get(query)
	if err != nil {
		fmt.Println(err)
	} else {
		contents, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println(err2)
		} else {
			json.Unmarshal(contents, &result)
		}
	}
	return fixScheduleOfBar(result.Bar)
}

func fixScheduleOfBar(bar GooglePlaceDetailed) GooglePlaceDetailed {
	for i := 0; i < len(bar.OpeningHours.Schedule); i++ {
		var twoParts = strings.SplitAfter(bar.OpeningHours.Schedule[i], "y")
		bar.OpeningHours.Schedule[i] = twoParts[1][2:len(twoParts[1])]
	}
	return bar
}

func getPlaceIDsOfBarsFromGooglePlacesAPI(latitude string, longitude string, radius string) []GooglePlace {
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?key=AIzaSyBDGTJegyakdJ3ObWRQfecI9zH_0MyzRhM&location=" + latitude + "," + longitude + "&radius=" + radius + "&type=bar")
	if err != nil {
		fmt.Println(err)
	} else {
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			var result Results
			json.Unmarshal(contents, &result)
			var bars = result.Bars[0:len(result.Bars)] // turning the array into a slice for easy use
			var resultBarsSlice []GooglePlace
			var i = 1
			for {
				bars = append(bars, resultBarsSlice...)
				if result.NextPageToken == "" {
					break
				}
				var query string = "https://maps.googleapis.com/maps/api/place/nearbysearch/json?key=AIzaSyBDGTJegyakdJ3ObWRQfecI9zH_0MyzRhM&pagetoken=" + result.NextPageToken
				time.Sleep(2 * time.Second)
				resp2, err2 := http.Get(query)
				if err2 != nil {
					fmt.Println("there was an error!")
				} else {
					contents2, err3 := ioutil.ReadAll(resp2.Body)
					if err3 != nil {
						fmt.Println(err3)
					} else {
						var resultForNextPage Results
						json.Unmarshal(contents2, &resultForNextPage)
						result = resultForNextPage
						resultBarsSlice = result.Bars[0:len(result.Bars)]
					}
				}
				i++
			}
			return bars
		}
	}
	return nil
}

func getBarsWithinThisSquareMileageFromTheseCoordinates(latitudeString string, longitudeString string, squareMileageString string) QueryResult {
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
		queryResult.Error = "getBarsWithinThisSquareMileageFromTheseCoordinates function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)

	latitude, latitudeErr := strconv.ParseFloat(latitudeString, 64)
	longitude, longitudeErr := strconv.ParseFloat(longitudeString, 64)
	squareMileage, squareMileageErr := strconv.Atoi(squareMileageString)
	if latitudeErr != nil {
		queryResult.Error = "getBarsWithinThisSquareMileageFromTheseCoordinates function: latitude parameter is messed up. " + latitudeErr.Error()
		return queryResult
	}
	if longitudeErr != nil {
		queryResult.Error = "getBarsWithinThisSquareMileageFromTheseCoordinates function: longitude parameter is messed up. " + longitudeErr.Error()
		return queryResult
	}
	if squareMileageErr != nil {
		queryResult.Error = "getBarsWithinThisSquareMileageFromTheseCoordinates function: squareMileage isn't an integer. " + squareMileageErr.Error()
		return queryResult
	}
	squareMileage += 4 // just in case
	// Finally
	var bars []BarData
	firstCall := true
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue

	for {
		var scanItemsInput = dynamodb.ScanInput{}
		scanItemsInput.SetTableName("Bar")
		if firstCall == false && lastEvaluatedKey == nil {
			break
		} else {
			scanItemsInput.SetExclusiveStartKey(lastEvaluatedKey)
		}
		expressionValuePlaceholders := make(map[string]*dynamodb.AttributeValue)
		// approximately (it's a square not a circle) getting bars within a 25 mile radius
		// 1 degree of latitude is approximately 69 miles
		// (1 degree of latitude / 69 miles) * (squareMileage / 2) = # degrees of latitude that equals half of the squareMileage
		degreesOfLatitudeWhichEqualHalfThisSquareMileage := (float64(1.00) / float64(69.00)) * (float64(squareMileage) / float64(2.00))
		// # of degrees / 25 miles = (1 degree of longitude / (69.1703 * COS(Latitude * 0.0174533)) ) * (squareMileage / 2)
		degreesOfLongitudeWhichEqualHalfThisSquareMileage := (1 / (69.1703 * math.Cos(latitude*0.0174533))) * (float64(squareMileage) / float64(2.00))
		latitudeSouth := latitude - degreesOfLatitudeWhichEqualHalfThisSquareMileage
		latitudeNorth := latitude + degreesOfLatitudeWhichEqualHalfThisSquareMileage
		longitudeEast := longitude - degreesOfLongitudeWhichEqualHalfThisSquareMileage
		longitudeWest := longitude + degreesOfLongitudeWhichEqualHalfThisSquareMileage
		latitudeSouthAttributeValue := dynamodb.AttributeValue{}
		latitudeNorthAttributeValue := dynamodb.AttributeValue{}
		longitudeEastAttributeValue := dynamodb.AttributeValue{}
		longitudeWestAttributeValue := dynamodb.AttributeValue{}
		latitudeSouthAttributeValue.SetN(strconv.FormatFloat(latitudeSouth, 'f', -1, 64))
		latitudeNorthAttributeValue.SetN(strconv.FormatFloat(latitudeNorth, 'f', -1, 64))
		longitudeEastAttributeValue.SetN(strconv.FormatFloat(longitudeEast, 'f', -1, 64))
		longitudeWestAttributeValue.SetN(strconv.FormatFloat(longitudeWest, 'f', -1, 64))
		expressionValuePlaceholders[":latitudeSouth"] = &latitudeSouthAttributeValue
		expressionValuePlaceholders[":latitudeNorth"] = &latitudeNorthAttributeValue
		expressionValuePlaceholders[":longitudeEast"] = &longitudeEastAttributeValue
		expressionValuePlaceholders[":longitudeWest"] = &longitudeWestAttributeValue
		scanItemsInput.SetExpressionAttributeValues(expressionValuePlaceholders)
		scanItemsInput.SetFilterExpression("(latitude BETWEEN :latitudeSouth AND :latitudeNorth) AND (longitude BETWEEN :longitudeEast AND :longitudeWest)")
		scanItemsOutput, err2 := getter.DynamoDB.Scan(&scanItemsInput)
		var dynamodbCall = DynamodbCall{}
		if err2 != nil {
			dynamodbCall.Error = "getBarsWithinThisSquareMileageFromTheseCoordinates function: Scan error. " + err2.Error()
			dynamodbCall.Succeeded = false
			queryResult.DynamodbCalls[0] = dynamodbCall
			queryResult.Error += dynamodbCall.Error
			return queryResult
		}
		dynamodbCall.Succeeded = true
		queryResult.DynamodbCalls[0] = dynamodbCall

		data := scanItemsOutput.Items
		barsOnThisPage := make([]BarData, len(data))
		jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &barsOnThisPage)
		if jsonErr != nil {
			queryResult.Error = "getBarsWithinThisSquareMileageFromTheseCoordinates function: UnmarshalListOfMaps error. " + jsonErr.Error()
			return queryResult
		}
		bars = append(bars, barsOnThisPage...)
		lastEvaluatedKey = scanItemsOutput.LastEvaluatedKey
		firstCall = false
	}

	queryResult.Bars = bars
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

func getAllBars() QueryResult {
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
		queryResult.Error = "getAllBars function: session creation error. " + err.Error()
		return queryResult
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var scanItemsInput = dynamodb.ScanInput{}
	scanItemsInput.SetTableName("Bar")
	scanItemsOutput, err2 := getter.DynamoDB.Scan(&scanItemsInput)
	var dynamodbCall = DynamodbCall{}
	if err2 != nil {
		dynamodbCall.Error = "getAllBars function: Scan error. " + err2.Error()
		dynamodbCall.Succeeded = false
		queryResult.DynamodbCalls[0] = dynamodbCall
		queryResult.Error += dynamodbCall.Error
		return queryResult
	}
	dynamodbCall.Succeeded = true
	queryResult.DynamodbCalls[0] = dynamodbCall
	data := scanItemsOutput.Items
	bars := make([]BarData, len(data))
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data, &bars)
	if jsonErr != nil {
		queryResult.Error = "getAllBars function: UnmarshalListOfMaps error. " + jsonErr.Error()
		return queryResult
	}
	queryResult.Bars = bars
	queryResult.DynamodbCalls = nil
	queryResult.Succeeded = true
	return queryResult
}

// Input = number of square miles of bars we want to populate the map with (1 mile = 1609.34 meters)
//
// an n x n matrix gives us an (n - 1)r x (n - 1)r square of bars, where r = 100 meters
// 		We need to determine what n's value should be.
//
//		Input should equal r * (n - 1)
// 				Input = numMiles = r * (n - 1)
// 							= (numMiles * 1609.34 meters / mile) = r meters * (n - 1)
//							=> n - 1 = (numMiles * 1609.34 meters) / (r meters * mile)
//							= n = ((numMiles * 1609.34) / r) + 1

func calculateWhatNShouldBe(numberOfSquareMilesString string, radiusString string) int {
	numberOfSquareMiles, err := strconv.ParseFloat(numberOfSquareMilesString, 64)
	radius, err2 := strconv.ParseFloat(radiusString, 64)
	if err != nil || err2 != nil {
		return 0
	}
	n := int((numberOfSquareMiles*1609.34)/radius) + 1
	if (n % 2) == 0 { // need n to be odd so that calulating lat and lng offset is easy
		n++
	}
	return n
}

type Results struct {
	NextPageToken string        `json:"next_page_token"`
	Bars          []GooglePlace `json:"results"`
}

type Result struct {
	Bar GooglePlaceDetailed `json:"result"`
}
