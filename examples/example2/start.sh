echo Killing all the running example processes
pkill stockweb 
pkill stockticker
pkill loginservice
pkill contentprovider
pkill database 
echo Wait 5 seconds before starting processes

sleep 5

cp ../../demoservice stockweb
cp ../../demoservice stockticker
cp ../../demoservice loginservice
cp ../../demoservice contentprovider
cp ../../demoservice database
echo copied all demo processes

./stockweb 9301 > stockweb.log 2>&1 &
./stockticker 9302 > stockticker.log 2>&1 &
./loginservice 9303 > loginservice.log 2>&1 &
./contentprovider 9304 > contentprovider.log 2>&1 &
./database 9305 > database.log 2>&1 &
echo started all demo services

sleep 5

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:9303", "Count" : 1 },
    { "Adr" : "http://localhost:9302", "Count" : 2 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:9301/config'

echo configured stockweb to call login service and stockticker

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:9303", "Count" : 2 },
    { "Adr" : "http://localhost:9304", "Count" : 2 }
  ]
}' \
 'http://localhost:9302/config'

echo configured stockticker to call authetication service and content service

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:9305", "Count" : 5 }
  ]
}' \
 'http://localhost:9303/config'

echo configured login service to call the database

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:9305", "Count" : 5 }
  ]
}' \
 'http://localhost:9304/config'

echo configured the content service to call the database