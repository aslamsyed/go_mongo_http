# mongo_http
GoLang Mongo DB REST API's

Test REST endpoints in REST client like Postman
Create a new Product
POST http://localhost:8080/product
{
    "id":    "100",
    "name":   "Amazon Echo",
    "title":  "Amazon Echo",
    "description":   "Amazon Echo is a best selling product"
    "price": "$299.99"
}

Get All products
GET http://localhost:8080/products

Get Product By Id
GET http://localhost:8080/product/100

Delete the Product
DELETE http://localhost:8080/product/100
