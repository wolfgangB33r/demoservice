curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "CrashConfig" : {
    "Code" : 9
  },
  "Callees" : [
  ]
}' \
 'http://localhost:8002/config'
