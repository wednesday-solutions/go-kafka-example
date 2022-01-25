#!/bin/bash
cd producer
go run cmd/server/main.go
cd ../consumer
go run cmd/server/main.go