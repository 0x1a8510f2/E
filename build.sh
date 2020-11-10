#!/bin/bash

ldflags="-X main.VcsCommit=$(git rev-parse --short HEAD) -X 'main.BuildTime=`date`'"

if [ $1 == "run" ]; then
    go run -ldflags "$ldflags" "."
else
    go build -ldflags "$ldflags" "."
fi