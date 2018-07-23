#### README

##### Project overview
Some information is  available from the server-side (Comandante) and not the client-side (Android)

Status quo: Tala's Android events were logged in Amplitude, but not server-side.

Project objective : write a AWS Lambda that reads server-side (Comandante) info from Kinesis and creates
                            a corresponding HTTP POST request to Amplitude.
                            
![Process Diagram](https://docs.aws.amazon.com/lambda/latest/dg/images/kinesis-pull-10.png)

##### Amplitude EVENT vs IDENTIFICATION
There are two types of posts -- event and identification.
1. Event posts are logged as events that can be viewed on a user's timeline. User properties sent with an event post will be reflected accordingly.
2. Identification posts can be viewed as "update user properties only" posts -- they are not sent with events; user properties added / changed with these posts will be reflected starting with the next Event post.  For example, consider the following stream of three consecutive posts for single user: 
 (_Event post 1_) -> (_Identification post 1: add property field "KYC status"_)  ->  (_Event post 2_) 
Event post 1 on the user's timeline will not have a "KYC status" field.
Identification post 1 will not appear on the user's timeline.
Event post 2 will reflect the "KYC status" field.

##### Design choices
* Used snake_case instead of CamelCase to reflect Tala's Comandante Style guide  .
Although Comandante is written in Kotlin (not Go), much of the Json information we handle with this lambda is formatted in Comandante; Lambda selects the relevant fields but otherwise doesn't do much to reformat.

##### What needs to happen on Comandante
Comandante must create appropriate json byte array. For rudimentary testing purposes, I added     the following lines to the main function in Comandante.  
    
``` 
val kinesisPublisher = KinesisPublisher()
kinesisPublisher.addUserRecord("""[{"user_id":"test-user-molly-9", "post_type": "event", "person_id": "111111123456", "event_type":"Cashout: Request", "country":"BRAZIL", "loan_application_id": "16"}] """.toByteArray())
```
Next steps would be to incrementally log the fields and build the string in Comandante.
   
##### What needs to happen on AWS
 Please refer to the AWS Kinesis stream **test-amplitude** and the AWS Lambda **test-amplitude-molly** in the __legacy__ environment, as well as the procedures on [this page](https://read.acloud.guru/serverless-golang-api-with-aws-lambda-34e442385a6a) .
 
* With this project, build executable "main" and zip main into deployment package "deployment.zip"
* Use execution role service-role/lambda-kinesis

**Important Notes**  
* Set GOOS to Linux
* Must set Handler name to the name of the executable (thus, in this case, set as "main")
