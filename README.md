# Generic demo service

Purpose of this service is to simulate specific anomaly situations, such as:

- Slowdowns
- Failures
- Increase in resource consumption
- Process crashes
- Calls to one or more other services

Therefore, this generic service can be used to build all kinds of variable demo 
situations with either flat or deep call trees.

The service listenes on following HTTP resource pathes:
- GET '/' normal page return
- POST '/config' service receives anomaly config as a HTTP POST message with a JSON config payload that is defined below.


## Config JSON body

{
  "error" : {
    
  },
  "slowdown" : {
    
  },
  "crash" : {
  
  },
  "cpu" : {
    
  },
  "memory" : {
    
  },
  "callees" : [
    { "adr" : "http://www.example.com", "count" : 10 },
    { "adr" : "http://www.orf.at", "count" : 3 }
  ]
}


