
# Go Mailer Service

## Pre-requisite
```
- aws cli

# configure aws cli
$ aws configure
    AWS Access Key ID: *******************
    AWS Secret Access Key: *******************
    Default region name: us-west-1
    Default output format: json

# then we also need to add new profile for aws credentials
$ vi ~/.aws/credentials
    ...

    [go-mailer-dev-account]
    aws_access_key_id = *******************
    aws_secret_access_key = *******************

    ...
```

## Development Build Setup
```
# go to src directory
$ cd go-mailer/src

# install dependencies
$ go install

# serve at localhost:9091
$ go run main.go

# build for production
$ go run build main.go
```

## Production Build Setup
```
# prerequisite
- docker
- docker-compose

# running the server
docker-compose up -d
```

## Running Unit Tests
```
# running the test
TODO: Add unit tests
```

## API Documentation

```
POST http://localhost:9091/api/send-email

Sample request:
Headers:
  Authorization: ad09f7f9-70c9-4a95-b1d1-3ef4f25a93f1
  Content-Type: "application/json"

Body:
  {
    "isMockEmail": true,
    "sender": "customer@gmail.com",
    "replyTo": "customer@gmail.com",
    "toRecipients": ["s.kelvinjohn@gmail.com", "kelvin@gmail.com"],
    "subject": "Amazon SES Test (AWS SDK for Go)",
    "htmlBody": "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with <a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the <a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>",
    "ccRecipients": [],
    "bccRecipients": [],
    "attachments": []
  }
```