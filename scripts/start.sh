#!/bin/bash

go build -o auth application/microservices/auth/cmd/main.go
go build -o sub application/microservices/subscription/cmd/main.go
go build -o chat application/microservices/chat/cmd/main.go
go build -o main cmd/main.go

nohup ./auth > auth.out &
nohup ./sub > sub.out &
nohup ./chat > chat.out &
nohup ./main > main.out &
