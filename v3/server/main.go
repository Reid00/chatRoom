/*
@Time        :2021/06/29 12:05:33
@Author      :Reid
@Version     :1.0
@Desc        : 服务器端的运行
*/

package main

func main() {
	server := newServer("127.0.0.1", 8888)
	server.Start()
}
