#!/bin/sh
go build -ldflags="-s -w" -o stripcontrol-golang cmd/service/main.go