<h2>Receipt Processor</h2>

### Assumptions:
1. The given purchaseDate is using the layout "2022-01-01".
2. The given purchaseTime is using the layout "13:01".
3. The given item prices and total are non negative.
4. The given item prices added up to total.


### Prerequisites

1. Docker
2. Golang

### API Specification
    Path: /receipts/process
    Method: POST
    Payload: Receipt JSON
    Response: JSON containing an id for the receipt.

    Path: /receipts/{id}/points
    Method: GET
    Response: A JSON object containing the number of points awarded.


### Getting Started

To start the app with Docker:
```
docker build -t my-app .
docker run -p 4000:4000 my-app
```

If Go is installed, to run without Docker:
```
go mod download
go run cmd/main.go
```


### Running the tests

Curl requests are provided in requests.txt, or you can run tests with the script
To run tests:
```
brew install jq
chmod +x test.sh
./test.sh
```