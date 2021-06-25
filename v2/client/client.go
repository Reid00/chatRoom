package main

import (
	"log"
	"net"
	"os"
)

// 客户端
func main() {

	// 创建跟服务器通信的socket
	conn, err := net.Dial("tcp", "127.0.0.1:8088")
	if err != nil {
		log.Fatalln("net.Dial error: ", err)
	}
	defer conn.Close()

	addr := conn.LocalAddr().String()
	sendMsg := "hello, this is " + addr
	conn.Write([]byte(sendMsg))

	// 获取用户键盘输入
	// fmt.scan 不能读取空格 回车
	// os.Stdin.Read() 可以，这里用这个
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				continue
			}
			// 把键盘输入发给服务器, 读多少，写多少

			// 如果客户端输入q、quit、exit等关闭客户端
			if string(buf[:n]) == "q\r\n" || string(buf[:n]) == "quit\r\n" || string(buf[:n]) == "exit\r\n" {
				os.Exit(1)
			} else {
				conn.Write([]byte(buf[:n]))
			}
		}
	}()

	// 从服务器读取消息
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			log.Println("server closed.")
			return
		}
		if err != nil {
			log.Println("conn.Read error: ", err)
			return
		}
		log.Println("server msg: ", string(buf[:n]))
	}
}
