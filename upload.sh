FILENAME="random"
FILESIZE=$(($RANDOM*100))
echo $(($FILESIZE/1024))kb
openssl rand -out $FILENAME $FILESIZE
curl --upload $FILENAME http://45.32.235.219:80/
rm "$FILENAME"
