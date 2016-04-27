FILENAME="random"

openssl rand -out $FILENAME $RANDOM
curl --upload $FILENAME http://localhost:8080/
rm "$FILENAME"
