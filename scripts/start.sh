#!/bin/bash

cd
go run ../application/microservices/auth/cmd/main.go
go run ../application/microservices/subscription/cmd/main.go
go run ../cmd/main.go