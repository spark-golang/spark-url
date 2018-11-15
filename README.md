# spark-url
Yet Another short url base on golang

# Base Use

## create a short url 

```shell
## crete a short url
curl -X "POST" "http://127.0.0.1:9999/local/url_create" \
     -H 'Content-Type: application/x-www-form-urlencoded; charset=utf-8' \
     --data-urlencode "url=https://golangcaff.com/docs/the-way-to-go/for-structure/37"
```

## request a short url 
```shell
url http://127.0.0.1:9999/Mw++  
```

## Use gRpc Demo 
1. gRpc server run 
```shell
 go run grpc/grpc/server.go
```

2. Run client 
```shell
go run main.go
```

3. Request 
```shell
curl -X "POST" "http://127.0.0.1:9999/hello" \
     -H 'Content-Type: application/x-www-form-urlencoded; charset=utf-8' \
     --data-urlencode "name=Panda"
```

# Thanks
Tim  
Gin   
Golang  
Google
