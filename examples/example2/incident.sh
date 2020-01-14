echo Triggering a slowdown of the database service by 50ms
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "SlowdownConfig" : {
    "SlowdownMillis" : 150,
    "Count" : 5000
  },
  "Callees" : [
  ]
}' \
 'http://localhost:9305/config'