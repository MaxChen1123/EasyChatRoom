# a Simple ChatRoom Implemented by Golang
to start the server:
```shell
go run ./Server/main.go -port 8888
```
to start clients:
```shell
go run ./Client/main.go -host 127.0.0.1 -port 8888 -name yourname
```

## Server
Server has a `broadcaster` to broadcast messages to every client

## Client
leave the ChatRoom by entering "Q"