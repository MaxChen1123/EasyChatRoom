package main

import(
    "fmt"
    "net"
    "bufio"
	"strings"
	"os"
    "flag"
)

var(
    host=flag.String("host","127.0.0.1","host")
    port=flag.String("port","8888","port")
    name=flag.String("name","user","name")
)

func serverReader(conn net.Conn){
    reader:=bufio.NewReader(conn)
    for{
        message,err:=reader.ReadString('\n')
        if err!=nil{
            fmt.Println("read from server---err:",err)
            return
        }
        fmt.Print(message)
    }
}

func main() {
    flag.Parse()
    conn, err := net.Dial("tcp", *host+":"+*port)
    if err != nil {
        fmt.Println("err :", err)
        return
    }
    defer conn.Close() 
    conn.Write([]byte(*name))
    fmt.Println("your name is:",*name)
    inputReader := bufio.NewReader(os.Stdin)
    go serverReader(conn)
    for {
        input, _ := inputReader.ReadString('\n') 
        inputInfo := strings.Trim(input, "\r\n")
        if strings.ToUpper(inputInfo) == "Q" { 
            return
        }
        _, err = conn.Write([]byte(inputInfo))
        if err != nil {
            fmt.Println("write to server---err:", err)
            return
        }
    }
}

