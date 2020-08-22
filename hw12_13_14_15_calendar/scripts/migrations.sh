#/bin/sh

file=tmp.sh
# Load all variables to tmp file
sed -e 's/:[^:\/\/]/=/g;s/$//g;s/ *=/=/g' configs/calendar.yml > $file
echo "goose -dir ./migrations postgres \"user=\$dbUser password=\$dbPass  dbname=\$dbName sslmode=disable\" up" >> $file
chmod +x $file
sh $file
rm $file