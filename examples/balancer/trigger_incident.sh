# first, configure the balancer to only send to 2 workers instead of previously 5

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "Callees" : [
    { "Adr" : "http://localhost:9400", "Count" : 3 },
    { "Adr" : "http://localhost:9500", "Count" : 3 }
  ],
  "Balanced" : true
}' \
 'http://localhost:9000/config'

# then, slowdown the remaining 2 workers because of the load shift

sleep 80

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "SlowdownConfig" : {
    "SlowdownMillis" : 500,
    "Count" : 3000
  },
  "Callees" : [
  ]
}' \
 'http://localhost:9400/config'

curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "SlowdownConfig" : {
    "SlowdownMillis" : 500,
    "Count" : 3000
  },
  "Callees" : [
  ]
}' \
 'http://localhost:9500/config'
