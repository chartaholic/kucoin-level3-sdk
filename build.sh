#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o kucoin-market-for-linux kucoin_market.go

CGO_ENABLED=0 go build -ldflags '-s -w' -o kucoin-market-for-mac kucoin_market.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o kucoin-market-for-windows.exe kucoin_market.go

