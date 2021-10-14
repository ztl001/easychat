package main

import(
	"fmt"
	"net"
)

type User struct{
	C     chan string
	Name  string
	Addr  string
}

var onlineMap = make(map[string]User)
var message   = make(chan string)
var isQuit    = make(chan bool)

func MakeMsg(cli User, msg string) (buf string){
	buf = "[" + cli.Addr + "]" + cli.Name + ": " + msg
	return
}

func HandleConn(conn net.Conn){
	defer conn.Close()

	cliAddr := conn.RemoteAddr().String()
	cli := User{make(chan string), cliAddr, cliAddr}
	onlineMap[cliAddr] = cli
	
	//广播信息
	go WriteMsgToClient(cli, conn)
	//message <- "[" + cli.Addr + "]" + cli.Name + ": loing"
	message <- MakeMsg(cli, "login")

	//接收用户信息
	go func(){
		buf := make([]byte, 2048)
		
		for{
			n, err := conn.Read(buf)
			//对方断开或出问题
			if n == 0{
				isQuit <- true

				fmt.Println("conn read err = ", err)
				return
			}

			msg := string(buf[:n-1])

			//转发消息
			message <- MakeMsg(cli, msg)
		}
	}()

	select{
		case <- isQuit:
			delete(onlineMap, cliAddr)
			message <- MakeMsg(cli, "login out")
			return
	}
}

func WriteMsgToClient(cli User, conn net.Conn){
	for msg := range cli.C{
		conn.Write([]byte(msg + "\n"))
	}
}

func Manager(){
	for {
		msg := <- message
		for _, cli := range onlineMap{
			cli.C <- msg
		}
	}
}


func main(){
	listener, err := net.Listen("tcp", ":8080")
	if err != nil{
		fmt.Println("net listen err = ", err)
		return
	}

	defer listener.Close()

	//转发消息
	go Manager()

	for{
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("listener accept err = ", err)
			continue
		}
		
		//处理用户链接
		go HandleConn(conn)
	}
}
