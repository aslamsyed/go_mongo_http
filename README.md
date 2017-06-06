# mongo_http
GoLang Mongo DB REST API's

mgo is a rich mongo driver for Go. 
http://labix.org/mgo

Execution Steps:

1) Setup mgo in local Go:
$ go get gopkg.in/mgo.v2
$ go get gopkg.in/mgo.v2/bson

2) Build the script:
$ go build

3) Run the script:
$ mongo_http

4) Test the REST endpoints in REST client like Postman
Create a new Product
Method: POST
URL: http://localhost:8080/product
{
    "id":    "100",
    "name":   "Amazon Echo",
    "title":  "Amazon Echo",
    "description":   "Amazon Echo is a best selling product"
    "price": "$299.99"
}

Get All products
Method: GET
URL: http://localhost:8080/products

Get Product By Id
Method: GET
URL: http://localhost:8080/product/100

Delete the Product
Method: DELETE
URL: http://localhost:8080/product/100
