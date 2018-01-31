# SMS Sender
SMS sending service with HTTP REST interface.

### Requirements
  * curl
  * docker
  * docker-compose

### Setup & Run
    docker-compose up

### API

#### Send SMS
    curl -XPOST http://localhost:8000/sms -d '{"from": "MY SERVICE", "to": "9271234567", "text": "Hello, world"}'

    {"status":"success","data":{"sms_id":"873487834"}}

#### Get SMS status
    curl -XGET http://localhost:8000/sms/873487834

    {"status":"success","data":{"sms_id":"873487834", "delivery_status": "delivered"}}
