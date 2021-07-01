/*
@Time        :2021/06/28 19:10:43
@Author      :Reid
@Version     :1.0
@Desc        :面向struct 的方式实现client端
*/
package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type client struct {
	name      string      // 客户端姓名
	msgChan   chan string // 服务器端返回客户端的消息channel
	inputChan chan string // 键盘输入给客户端的消息channel
}

// 构造函数
func newClient(conn net.Conn) *client {
	return &client{
		name:      conn.LocalAddr().String(),
		msgChan:   make(chan string),
		inputChan: make(chan string),
	}
}

// 获取键盘输入
func (c *client) getInput() {
	buf := make([]byte, 4096)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			continue
		}
		// 读多少, 写多少
		if input := strings.ToLower(strings.TrimSpace(string(buf[:n]))); input == "q" || input == "quit" || input == "exit" {
			os.Exit(1)
		} else {
			input := strings.TrimSpace(string(buf[:n]))
			c.inputChan <- input
		}
	}
}

// 从服务器读取数据
func (c *client) readServer(conn net.Conn) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			log.Println("Server closed, client to close...")
			return
		}
		if err != nil {
			log.Println("conn.Read error: ", err)
			return
		}
		msg := strings.TrimSpace(string(buf[:n]))
		c.msgChan <- msg
	}

}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("net.Dial error: ", err)
	}
	defer conn.Close()
	addr := conn.LocalAddr().String()
	login := "hello, this is " + "[" + addr + "]"
	conn.Write([]byte(login))

	c := newClient(conn)
	go c.getInput()
	go c.readServer(conn)

	for {
		select {
		case msg := <-c.msgChan:
			log.Println("server message: ", msg)
		case input := <-c.inputChan:
			conn.Write([]byte(input))
		case <-time.After(time.Second * 10):
			log.Println("十秒内没有通信，退出")
			return
		}
	}

}
