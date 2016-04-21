openssl rand -out abc.db $RANDOM
curl -F "upload=@abc.db" http://localhost:8080/
curl --upload abc.db http://localhost:8080/
rm "abc.db"
