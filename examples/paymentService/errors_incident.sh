curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "ErrorConfig" : {
    "ResponseCode" : 500,
    "Count": 100
  },
  "Callees" : [
  ]
}' \
 'http://localhost:8083/config'