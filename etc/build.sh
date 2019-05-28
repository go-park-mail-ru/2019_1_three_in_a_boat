#!/usr/bin/env bash
protoc -I shared/formats/proto shared/formats/proto/auth.proto --go_out=plugins=grpc:shared/formats/pb
go build -o run-auth ./auth
go build -o run-server ./server
go build -o run-game ./game
go build -o run-chat ./chat
