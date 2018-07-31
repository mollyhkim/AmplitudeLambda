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


/* Retry behavior: 
 * When a Lambda function invocation fails, AWS Lambda attempts to process the erring batch 
 * of records until the time the data expires, which can be up to seven days.
 * The exception is treated as blocking, and AWS Lambda will not read any new records from the shard 
 * until the failed batch of records either expires or is processed successfully. This ensures that 
 * AWS Lambda processes the stream events in order.
 *
 */