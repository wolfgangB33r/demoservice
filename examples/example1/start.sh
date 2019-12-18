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

./customerAPI 8080 >log.log&
./sideService 8081&
./calcService 8082&
./restService 8083&
./recomService 8084&
./modelService 8085&

sleep 5

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8081", "Count" : 2 },
    { "Adr" : "http://localhost:8082", "Count" : 3 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8080/config'

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8083", "Count" : 2 },
    { "Adr" : "http://localhost:8084", "Count" : 2 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8081/config'

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8085", "Count" : 2 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8084/config'
