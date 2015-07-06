#!/bin/sh
if [ "$CI" = "true" ]
then
    npm update -g npm
    npm install -g dredd
fi
go build main.go
./main -m &
sleep 3
PID=$!
dredd apiary.apib http://localhost:8080/
RESULT=$?
kill -9 $PID
exit $RESULT
