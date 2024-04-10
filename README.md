# Sharer

## Easy way to share files without need a middle man

### Not fully baked

### Advantages

* No Middle Man

* No share limit, 20 MB by default that can be modified

* The information is transmitted in binary so it is very fast

* Directly from the CLI which is quite convenient
---
### Sharer [ client | server ]

## Flags
## share server
  - ### ( -port -p ) port to listen (default 9000). 
  - ### ( -size | -s ) size (MB) allow to recieve per request (default 20MB).
  - ### ( -token | -t ) token to allow client share to you (optional).
  - ### ( -unsecure | -u ) to not use token and allow everyone share to you. 
        sharer server -p 4000 -s 30 -t mycustomtoken

## share client
  - ### ( -file | -f ) path to file|folder  
  - ### ( -addr | -a ) set address to send\n  
  - ### ( -token | -t ) set token to authenticate with server (optional) 
        sharer client -f myfile.txt -a 127.0.0.1:9000 -t myauthtoken

