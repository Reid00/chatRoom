/*
@Time        :2021/06/25 14:12:13
@Author      :Reid
@Version     :1.0
@Desc        :服务端，时刻监听客户端连接
*/

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// 创建用户结构体
type User struct {
	name    string      // 用户名称
	addr    net.Addr    // 用户ip地址
	msgChan chan string // 用户消息channel
}

func newUser(name string, addr net.Addr, msgChan chan string) *User {
	return &User{
		name:    name,
		addr:    addr,
		msgChan: msgChan,
	}
}

// 用户在线列表 [key] addr, [val] User 指针
var onlineMap = make(map[string]*User)

// 定义消息管理中心，由此广播给每个用户消息
var Message = make(chan string, 1)

// 消息广播中心，所有消息由此分发
func Manager() {
	for {
		select {
		// 读到消息后，发送给每个在线的user
		case msg := <-Message:
			//　广播消息 给每个在线User
			for _, user := range onlineMap {
				user.msgChan <- msg
				// 每个user 返回给各自的客户端
				// 目的: 消息广播给每个客户端， 如果去做，服务器写给每个客户端
				// 做法: 服务器一直读取对应User 的msgChan
			}
		default:
		}
	}
}

// 处理处理客户端连接请求，用于通信的socket
func HandleConn(conn net.Conn) {
	defer conn.Close()

	// 获取客户端地址
	addr := conn.RemoteAddr().String()
	// 用户信息
	user := newUser(addr, conn.RemoteAddr(), make(chan string))
	// 添加到在线列表中
	onlineMap[addr] = user

	msg := "欢迎用户: " + user.name
	Message <- msg

	// 获取键盘输入
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if n == 0 {
				log.Println("服务器被强制关闭")
				continue
			}
			if err != nil {
				log.Println("获取input失败: ", err)
				continue
			}
			if string(buf[:n]) == "q\r\n" || string(buf[:n]) == "quit\r\n" || string(buf[:n]) == "exit\r\n" {
				log.Println("收到服务器下线请求，关闭服务器")
				conn.Write([]byte("收到服务器下线请求，关闭服务器"))
				os.Exit(0)
			} else if string(buf[:n]) == "show\r\n" {
				// 遍历在线User 打印address
			}
		}
	}()

	// 从自己的消息列表钟读取信息
	go func() {
		for {
			// 读取自己的msgChan
			select {
			case msg := <-user.msgChan:
				conn.Write([]byte(msg))
			default:
			}
		}
	}()
	buf := make([]byte, 4096)
	// 读取客户端发来的消息 写到 Message
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			msg := fmt.Sprintf("[%s] 退出群聊...", addr)
			Message <- msg
			return
		}
		if err != nil {
			msg := fmt.Sprintf("从%s的客户端conn.Read读取失败.", addr)
			Message <- msg
			continue
		}
		msg := fmt.Sprintf("[%s]的消息: %s", addr, strings.TrimSpace(string(buf[:n])))
		Message <- msg
	}

}

func main() {
	// 创建监听连接的socket
	listener, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		log.Fatalln("net.Listen error: ", err)
	}
	defer listener.Close()

	// 用户管理， 消息广播
	go Manager()

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
