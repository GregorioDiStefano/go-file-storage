#!/bin/bash


#Run complied controller test from main path
(cd controllers && go test -c && cp controllers.test ../)
./controllers.test
rm controllers.test
