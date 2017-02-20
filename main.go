package main

import (
	"encoding/json"
	//"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/dynamoDB", queryDynamo)

	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func queryDynamo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Querying DynamoDB now."))
	data, err := queryPersonInDynamoDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}


func queryPersonInDynamoDB() (string, error) {
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2")}))
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
	    return "", err
	}

	var tables string = ""
	for _, table := range result.TableNames {
	    tables = tables + *table
	}
	return tables, nil
	/*
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config *aws.Config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("err")
	}
	var svc *dynamodb.DynamoDB = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var getItemInput = dynamodb.GetItemInput{}
	getItemInput.SetTableName("Person")
	var attributeValue = dynamodb.AttributeValue{}
	attributeValue.SetS("12345699033")
	getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"facebookID": &attributeValue})
	getItemOutput, err2 := getter.DynamoDB.GetItem(&getItemInput)
	return getItemOutput.GoString(), err2
	*/
}

type personData struct {
	BlackoutList struct {
		NS []string `json:"NS"`
	} `json:"blackoutList"`
	FacebookID struct {
		S string `json:"S"`
	} `json:"facebookID"`
	InvitedTo struct {
		NS []int64 `json:"NS"`
	} `json:"invitedTo"`
	IsABarOwner struct {
		BOOL bool `json:"BOOL"`
	} `json:"isABarOwner"`
	IsMale struct {
		BOOL bool `json:"BOOL"`
	} `json:"isMale"`
	Name struct {
		S string `json:"S"`
	} `json:"name"`
	PartyHostFor struct {
		NS []int64 `json:"NS"`
	} `json:"partyHostFor"`
}

/*
func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=YOUR_API_KEY&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}
*/

type partyData struct {
	AddressLine1        string  `json:"addressLine1"`
	AddressLine2        string  `json:"addressLine2"`
	City                string  `json:"city"`
	Country             string  `json:"country"`
	Description         string  `json:"description"`
	DrinksProvided      bool    `json:"drinksProvided"`
	EndTime             string  `json:"endTime"`
	FeeForDrinks        bool    `json:"feeForDrinks"`
	Latitude            float32 `json:"latitude"`
	Longitude           float32 `json:"longitude"`
	PartyID             int64   `json:"partyID"`
	PartyName           string  `json:"partyName"`
	StartTime           string  `json:"startTime"`
	StateProvinceRegion string  `json:"stateProvinceRegion"`
	ZipCode             int16   `json:"zipCode"`
}
