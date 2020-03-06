
# Go Mailer Service

## Development Build Setup
``` bash
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
``` bash
# prerequisite
- docker
- docker-compose

# running the server
docker-compose up -d
```

## Running Unit Tests
``` bash
# running the test
TODO: Add unit tests
```

## API Documentation

``` bash
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