# portfolio_svc

This repositry is a written in Golang 

The Portfolio Backend service is responsible for handling all the requests made by **portfolio_ui** service. 

## How to Run
`go run cmd/main/main.go`

Navigatge to `localhost:5050/healthy` for health check of the server

Repository Layout is based on golang community recomneded best practices. More on it [here](https://github.com/golang-standards/project-layout) 

## Libraries 

- gin is a highly scalable, light weight http server. 
- gorm to connect to a relational database 

## Adding a new ENV variable
1. add it to app.env
1. add the varaible to struct in `configs/env-config.go`
