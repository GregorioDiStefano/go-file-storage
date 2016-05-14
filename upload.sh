FILENAME="random"

openssl rand -out $FILENAME $(($RANDOM*1000))
curl --upload $FILENAME http://localhost:8080/
rm "$FILENAME"
