# financial-transactions-go-webapp

## Objective
>Web application to analyse and investigate financial transactions. 
Developed using Go Lang. 
Project created by AluraChallenge Backend 3 #alurachallengebackend3

## Requeriments

1) Go lang >= 1.17.7
2) Docker
3) Docker compose

## Development setup
1) Clone project from github
2) Download project dependencies:
```
go mod download
```
3) Start postgresql from a docker image:
```
docker compose up
```

## Execute Project
```
go run main.go
```
When the project starts, open it from localhost:8080/

*** The example spredsheet and xml transaction file are in folder /upload

## Features
- `Access Control`: Users CRUD with Login, Logout, Registration and Deletion of accounts
- `File upload`: CSV and XML files with financial transactions data to be analysed
- `Database Storage`: SQL Database Persistence
- `Transactions Analysis`: Investigation for fraudulent or suspicious transactions
