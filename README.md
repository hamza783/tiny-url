# tiny-url

This is a microservices project to create tiny url. Written in golang. It uses REST api to connect to api gateway which then uses gRPC to communicate to other microservices. The data is then saved in redis.
This is a small project to learn about implementing microservices arcitecture. All the services are saved in 1 repo for easy management.


## Setup
### api-gateway
To start the api gateway server: go run api-gateway/cmd/main.go
Uses port: 8080

### url-shortening-service
To start the service: go run url-shortening-service/cmd/main.go
Uses port: 8081

### url-redirection-service
To start the service: go run url-redirection-service/cmd/main.go
Uses port: 9000

### shared
Contains code that is shared between different services

### redis
make sure you have redis installed locally. Redis is used to store the short_url and long_url mapping. It is acting as a database for this project.
To start redis: redis-server

## Testing
### To shorten a url
POST `localhost:8080/api/shorten`
Request Example: 
{
    "long_url": "google.com"
}
Response Example:
{
    "data": {
        "long_url": "facebook.com",
        "short_url": "nebJrn"
    }
}

### To redirect to the long url using the short url
GET `localhost:8080/api/{short_url}`
Example url: localhost:8080/api/nebJrn
