#!/bin/bash

rm -f out/ustb-login*
cd src/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../out/ustb-login.linux.amd64 main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../out/ustb-login.darwin main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../out/ustb-login.win64.exe main.go
cd ..
