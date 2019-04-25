#!/usr/bin/env bash
protoc -I formats/proto formats/proto/auth.proto --go_out=plugins=grpc:formats/pb
go build -o run-auth ./auth
go build -o run-server ./server
