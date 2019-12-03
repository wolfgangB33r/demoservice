pkill customerAPI 
pkill sideService 
pkill calcService 
pkill restService 
pkill recomService 
pkill modelService 

sleep 5

cp ../../demoservice customerAPI 
cp ../../demoservice sideService 
cp ../../demoservice calcService 
cp ../../demoservice restService 
cp ../../demoservice recomService 
cp ../../demoservice modelService 

./customerAPI 8001&
./sideService 8002&
./calcService 8003&
./restService 8004&
./recomService 8005&
./modelService 8006&

sleep 5

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8002", "Count" : 2 },
    { "Adr" : "http://localhost:8003", "Count" : 6 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8001/config'

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8004", "Count" : 2 },
    { "Adr" : "http://localhost:8005", "Count" : 6 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8003/config'

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8006", "Count" : 2 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8005/config'
