package main

import (
	"fmt"
	"net"
	"time"
    "flag"
    "server/cli"
)

var (
    port =flag.String("port","8888","the port number")
)

type Client =*cli.ClientModel

var (
    messages = make(chan string)
    leaving =make(chan Client)
    entering =make(chan Client)
)

var clients = make(map[Client]bool)

func broadcaster(){
    for{
        select{
        case msg:=<-messages:
            for cli :=range clients{
                cli.Channel<-msg
            }
        case newCli:=<-entering:
            clients[newCli]=true
        case delCli:=<-leaving:
            delete(clients,delCli)
            delCli.Close()
        }
    }
}

func process(conn net.Conn) {
    var nameByte [128]byte 
    n,err:=conn.Read(nameByte[:])
    if err!=nil{
        fmt.Printf("something wrong:%s\n",err)
        return
    }
    who:=string(nameByte[:n])
    client:=cli.NewClientModel(conn,who)
    go client.DetectTimeout(leaving)

    entering<-client
    messages<-who+" has connected at "+getTimeNow()

    go client.WriteMsg()
    fmt.Println(who,"has connected")

    for {
        msg,err:=client.ReadMsg()
        if err!=nil{
            fmt.Printf(client.Name+" has something wrong:%s \n",err)
            leaving<-client
            messages<-client.Name+" has leaved at "+getTimeNow()
            return
        }
        messages<-client.Name+":"+msg+"\t\t"+getTimeNow()
    }
}

func getTimeNow() string {
    return time.Now().Format("15:04:05")
}

func main() {
    flag.Parse()
    listen, err := net.Listen("tcp", "127.0.0.1:"+*port)
    if err != nil {
        fmt.Println("listen failed, err:", err)
        return
    }
    go broadcaster()
    defer listen.Close()
    for {
        conn, err := listen.Accept()
        if err != nil {
            fmt.Println("accept failed, err:", err)
            continue
        }
        go process(conn)
    }
}