# SMS Sender
SMS sending service with HTTP REST interface (using https://bsms.tele2.ru/api).

### Requirements
  * curl
  * docker
  * docker-compose

### Setup & Run
    docker-compose up

### API

#### Send SMS
    curl -XPOST http://localhost:8000/sms -d '{"phone": "9271234567", "text": "Hello, world"}'

    {"status":"success","data":{"sms_id":"873487834"}}
    {"status":"error","data":{"error":"message"}}

#### Get SMS status
    curl -XGET http://localhost:8000/sms/873487834

    {"status":"success","data":{"delivery_status": "delivered"}}
    {"status":"error","data":{"error":"message"}}
