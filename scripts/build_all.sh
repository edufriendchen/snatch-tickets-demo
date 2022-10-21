#!/bin/bash

GOEXEC=${GOEXEC:-"go"}

# clean
rm -rf output/ && mkdir -p output/bin/ && mkdir -p output/log/

$GOEXEC mod tidy && go mod verify

# build clients
$GOEXEC build -v -o output/bin/client ./client/client.go

# build servers
$GOEXEC build -v -o output/bin/ticketing_server ./server/ticketing_server.go
$GOEXEC build -v -o output/bin/order_server ./server/order_server.go
