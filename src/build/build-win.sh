#!/bin/bash

CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build  --trimpath -o ../../bin/mudgateway.exe ../