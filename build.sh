#!/usr/bin/env bash
cd auth
protoc -I proto proto/auth.proto --go_out=plugins=grpc:proto
go build -o ../run-auth .
cd ..
go build -o run-server