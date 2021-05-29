#!/bin/bash

go build -o auth_service application/microservices/auth/cmd/main.go
go build -o sub_service application/microservices/subscription/cmd/main.go
go build -o chat_service application/microservices/chat/cmd/main.go
go build -o kudago_service application/microservices/api_kudago/cmd/main.go
go build -o main_service cmd/main.go

nohup ./auth_service > auth.out &
nohup ./sub_service > sub.out &
nohup ./chat_service > chat.out &
nohup ./kudago_service > kudago.out &
nohup ./main_service > main.out &
