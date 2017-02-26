package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Find all parties I'm invited to
func myParties(w http.ResponseWriter, r *http.Request) {
	// http://example.com/page?parameter=value&also=another
	params := strings.Split(r.URL.Query().Get("partyID"), ",")
	data, err := findMyParties(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	parties := make([]PartyData, 1)
	jsonErr := dynamodbattribute.UnmarshalListOfMaps(data["Party"], &parties)

	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(parties)
}

func findMyParties(params []string) (map[string][]map[string]*dynamodb.AttributeValue, error) {
	type ItemGetter struct {
		DynamoDB dynamodbiface.DynamoDBAPI
	}
	// Setup
	var getter = new(ItemGetter)
	var config = &aws.Config{Region: aws.String("us-west-2")}
	sess, err := session.NewSession(config)
	if err != nil {
		fmt.Println("err")
	}
	var svc = dynamodb.New(sess)
	getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)
	// Finally
	var batchGetItemInput = dynamodb.BatchGetItemInput{}
	attributesAndValues := make([]map[string]*dynamodb.AttributeValue, len(params))
	for i := 0; i < len(params); i++ {
		var attributeValue = dynamodb.AttributeValue{}
		attributeValue.SetN(params[i])
		attributesAndValues[i] = make(map[string]*dynamodb.AttributeValue)
		attributesAndValues[i]["partyID"] = &attributeValue
	}
	/*
		var attributeValue1 = dynamodb.AttributeValue{}
		var attributeValue2 = dynamodb.AttributeValue{}
		attributeValue1.SetN("1")
		attributeValue2.SetN("2")
		attributesAndValues[0] = make(map[string]*dynamodb.AttributeValue)
		attributesAndValues[1] = make(map[string]*dynamodb.AttributeValue)
		attributesAndValues[0]["partyID"] = &attributeValue1
		attributesAndValues[1]["partyID"] = &attributeValue2
	*/

	var keysAndAttributes dynamodb.KeysAndAttributes
	keysAndAttributes.SetKeys(attributesAndValues)

	requestedItems := make(map[string]*dynamodb.KeysAndAttributes)
	requestedItems["Party"] = &keysAndAttributes
	batchGetItemInput.SetRequestItems(requestedItems)
	//getItemInput.SetKey(map[string]*dynamodb.AttributeValue{"partyID": &attributeValue})
	batchGetItemOutput, err2 := getter.DynamoDB.BatchGetItem(&batchGetItemInput)
	return batchGetItemOutput.Responses, err2
}

/*
{
    "TableName" : {
      keys: [
        {
          "AttributeName" => "value", # value <Hash,Array,String,Numeric,Boolean,IO,Set,nil>
        },
      ],
    }
}*/
