/*
@Time        :2021/06/29 11:13:29
@Author      :Reid
@Version     :1.0
@Desc        :群聊通信服务器端功能
*/
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// 定义一个server 对象
type Server struct {
	Ip   string
	Port int
}

// 定义一个server 构造函数
func newServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

// 处理业务请求
func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()
	log.Println("连接创建成功...")
	// 读取客户端数据
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			log.Printf("[server]: client %v closed.", conn.RemoteAddr())
			return
		}
		if err != nil {
			log.Println("[server]: conn.Read err: ", err)
			return
		}
		log.Println(string(buf[:n]))
	}
	io.Copy()
}

// 启动server方法
func (s *Server) Start() {
	// 创建监听连接的socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Fatal("net.Listen error: ", err)
	}
	// 关闭listener
	defer listener.Close()
	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("listener.Accept error: ", err)
		}

		// 处理请求
		go s.Handle(conn)

	}
}

// 读取客户端发送的消息
func (s *Server) ReadFromClient(conn net.Conn) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			log.Println("[客户端: " + conn.RemoteAddr().String() + "]" + "断开连接")
			return
		}
		if err != nil {
			log.Println("conn.Read error: ", err)
			return
		}
		// 处理读的数据
		if input := strings.ToLower(strings.TrimSpace(string(buf[:n]))); len(input) > 0 {
			conn.Write([]byte(strings.ToUpper(input)))
		}
	}
}
