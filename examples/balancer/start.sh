pkill balancer 
pkill worker 

sleep 5

cp ../../demoservice balancer 
cp ../../demoservice worker 

./balancer 9000&
./worker 9100&
./worker 9200&
./worker 9300&
./worker 9400&
./worker 9500&

sleep 5

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:9100", "Count" : 1 },
    { "Adr" : "http://localhost:9200", "Count" : 1 },
    { "Adr" : "http://localhost:9300", "Count" : 1 },
    { "Adr" : "http://localhost:9400", "Count" : 1 },
    { "Adr" : "http://localhost:9500", "Count" : 1 }
  ]
}' \
 'http://localhost:9000/config'
