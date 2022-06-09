curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "ResourceConfig" : {
    "Severity" : 100,
    "Count" : 30000
  },
  "Callees" : [
  ]
}' \
 'http://localhost:8083/config'
