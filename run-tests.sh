#!/bin/bash
UPLOAD_TO="http://localhost:8080/"

function upload_file {
    local FILENAME="random"
    openssl rand -out $FILENAME $RANDOM
    LINKS=$(curl --upload $FILENAME $UPLOAD_TO)
    rm "$FILENAME"
    #echo $LINKS | grep -Po '"downloadURL":.*?[^\\]",'
}


#Run complied controller test from main path
go build .

#upload_file

(cd controllers && go test -cover -c && cp controllers.test ../)
./controllers.test uploadRoutes_test.go
rm controllers.test
