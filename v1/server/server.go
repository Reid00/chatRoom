package main

import (
	"log"
	"net"
	"strings"
)

// 服务端，时刻监听客户端连接

// 处理客户端发送过来的请求
func HandleMsg() {

}

// 处理通信的socket
func HandleConn(conn net.Conn) {
	defer conn.Close()

	// 获取客户端地址
	addr := conn.RemoteAddr().String()
	log.Println("欢迎用户: ", addr)

	// 读取客户端发来的消息
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			log.Printf("%s 客户端关闭连接.", addr)
			return
		}
		if err != nil {
			log.Println("conn.Read error: ", err)
			return
		}
		// 读多少，处理多少用
		log.Printf("来自%s的信息为: %s", addr, string(buf[:n]))

		// 消息处理后， 返回给客户端消息
		msg := strings.Replace(string(buf[:n]), "?", "", -1)
		msg = strings.Replace(msg, "？", "", -1)
		msg = strings.Replace(msg, "吗", "", -1)
		//　去掉client 发来的\r\n
		msg = strings.Replace(msg, "\r\n", "", -1)
		msg = strings.Replace(msg, "\n", "", -1)
		msg = strings.ToUpper(msg)
		conn.Write([]byte(msg))
	}
}

func main() {
	// 创建监听连接的socket
	listener, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		log.Fatalln("net.Listen error: ", err)
	}
	defer listener.Close()

	log.Println("监听等待客户端连接...")
	// 循环等待客户端连接，可以连接多个客户端
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error: ", err)
		}

		// 创建一个goroutine 处理通信
		go HandleConn(conn)
	}
}
