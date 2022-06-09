pkill paymentService 
pkill authenticationService
pkill persistanceService
pkill riskAssessmentService

sleep 5

cp ../../demoservice paymentService
cp ../../demoservice authenticationService
cp ../../demoservice persistanceService
cp ../../demoservice riskAssessmentService

./paymentService 8080 > payment.log 2>&1 &
./authenticationService 8081 > payment.log 2>&1 &
./persistanceService 8082 > payment.log 2>&1 &
./riskAssessmentService 8083 > payment.log 2>&1 &

sleep 5

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:8081", "Count" : 1 },
    { "Adr" : "http://localhost:8082", "Count" : 5 },
    { "Adr" : "http://localhost:8083", "Count" : 1 }
  ]
}' \
 'http://localhost:8080/config'

