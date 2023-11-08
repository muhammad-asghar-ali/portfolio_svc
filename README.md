# portfolio_svc

This repositry is a written in Golang 

The Portfolio Backend service is responsible for handling all the requests made by **portfolio_ui** service. 

## How to Run
`go run cmd/main/main.go`

Navigatge to `localhost:8080/healthy` for health check of the server

Repository Layout is based on golang community recomneded best practices. More on it [here](https://github.com/golang-standards/project-layout) 

## Libraries 

- gin is a highly scalable, light weight http server. 
- gorm to connect to a relational database 
