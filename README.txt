I. Project overview
    Project written in Golang to write a AWS Lambda that reads server-side (Comandante) info from Kinesis and creates
    a corresponding HTTP POST request to Amplitude

II. Amplitude EVENT vs IDENTIFICATION
    There are two types of posts -- event and identification.
        A. Event posts are logged as events that can be viewed on a user's timeline. User properties sent with an
            event post will be reflected accordingly.
        B. Identification posts can be viewed as "update user properties only" posts -- they are not sent with events;
            user properties added / changed with these posts will be reflected starting with the next Event post
            e.g., Consider the following stream of three consecutive posts for single user:
               (Event post 1)  ->  (Identification post 1: add property field "KYC status")  ->  (Event post 2)

               Event post 1 on the user's timeline will not have a "KYC status" field.
               Identification post 1 will not appear on the user's timeline.
               Event post 2 will reflect the "KYC status" field.

III.Design choices
    A. Used snake_case instead of CamelCase to reflect Tala's Comandante Style guide
       Although Comandante is written in Kotlin (not Go), much of the Json information we handle with this lambda
       is formatted in Comandante; Lambda selects the relevant fields but otherwise doesn't do much to reformat.

        i. Constants in UPPER_SNAKE_CASE
