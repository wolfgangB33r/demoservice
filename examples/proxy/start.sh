echo Killing all the running example processes
pkill caller 
pkill proxy
pkill destination
echo Wait 5 seconds before starting processes

sleep 5

cp ../../demoservice caller
cp ../../demoservice proxy
cp ../../demoservice destination
echo copied all demo processes

./caller 8497 &
./proxy 8498 > proxy.log 2>&1 &
./destination 8499 &
echo started all demo services

sleep 5

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8498", "Count" : 1 }
  ]
}' \
 'http://localhost:8497/config'

echo configured caller to call the proxy

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8499", "Count" : 1 }
  ],
  "Proxy" : true
}' \
 'http://localhost:8498/config'

echo configured the proxy to call the destination service 
