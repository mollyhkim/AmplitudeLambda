package main

import (
	"fmt"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func kinesisHandler(kinesisEvent events.KinesisEvent) error {
	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data
		dataText := string(dataBytes)
		fmt.Printf("%s Data = %s \n", record.EventName, dataText)

		parsed := parsingInputFromCMD(dataText)
		switch parsed.Post_type {
		case "event":
			createEventJson(parsed)
		case "identification":
			createPropertiesJson(parsed)
		default:
			log.Fatal(`Input json post type must be either "event" or "identification"; received: ` + parsed.Post_type)
		}
	}

	return nil
}

func main() {
	fmt.Println("starting (i.e., compilation successful)")
	lambda.Start(kinesisHandler)
}
