sleep 5

cp ../../demoservice customerAPI
cp ../../demoservice sideService
cp ../../demoservice calcService
cp ../../demoservice restService
cp ../../demoservice recomService
cp ../../demoservice modelService

./customerAPI 8080&
./sideService 8082&
./calcService 8083&
./restService 8084&
./recomService 8085&
./modelService 8086&

sleep 5

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8082", "Count" : 2 },
    { "Adr" : "http://localhost:8083", "Count" : 3 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8080/config'

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8084", "Count" : 2 },
    { "Adr" : "http://localhost:8085", "Count" : 2 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8084/config'

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8086", "Count" : 2 },
    { "Adr" : "http://www.example.com", "Count" : 1 }
  ]
}' \
 'http://localhost:8085/config'
