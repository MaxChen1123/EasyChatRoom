package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
    "flag"
)

var (
    port =flag.String("port","8888","the port number")
)
type client chan string
var (
    messages = make(chan string)
    leaving =make(chan client)
    entering =make(chan client)
)

func broadcaster(){
    clients:=make(map[client]bool)
    for{
        select{
        case msg:=<-messages:
            for cli :=range clients{
                cli<-msg
            }
        case newCli:=<-entering:
            clients[newCli]=true
        case delCli:=<-leaving:
            delete(clients,delCli)
            close(delCli)
        }
    }
}

func clientWriter(conn net.Conn,ch client){
    for msg:=range ch{
        fmt.Fprintln(conn,msg)
    }
}

func process(conn net.Conn) {
    defer conn.Close() 
    ch := make(client)
    var nameByte [128]byte 
    n,err:=conn.Read(nameByte[:])
    if err!=nil{
        fmt.Printf("something wrong:%s\n",err)
        return
    }
    who:=string(nameByte[:n])
    entering<-ch
    messages<-who+" has connected at "+getTimeNow()
    reader := bufio.NewReader(conn)
    go clientWriter(conn,ch)
    fmt.Println(who,"has connected")
    var buf [128]byte
    for {
        n,err:=reader.Read(buf[:])
        if err!=nil{
            fmt.Printf("something wrong:%s\n",err)
            leaving<-ch
            messages<-who+" has leaved at "+getTimeNow()
            return
        }
        recvStr := string(buf[:n])
        messages<-who+":"+recvStr+"\t\t"+getTimeNow()
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