package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// 1. 准备 TCP 地址 (去掉 http://，并指定端口号)
	// 如果你的服务在其他端口，请把 80 改为对应的端口，例如 "gw.fps.dev.zrzkwlw.com:9000"
	address := "gw.fps.dev.zrzkwlw.com:80"

	// 2. 准备要发送的十六进制报文
	hexStr := "40405E0001010B0E1014051A86EB242C2F00000000000000070002CA010000000000802323"

	// 将十六进制字符串解码为 byte 数组
	payload, err := hex.DecodeString(hexStr)
	if err != nil {
		log.Fatalf("报文解码失败: %v", err)
	}

	// 3. 建立 TCP 连接
	fmt.Printf("正在连接到 %s...\n", address)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second) // 5秒超时
	if err != nil {
		log.Fatalf("TCP 连接失败: %v", err)
	}
	defer conn.Close()
	fmt.Println("连接成功！")

	// 4. 发送报文
	n, err := conn.Write(payload)
	if err != nil {
		log.Fatalf("发送报文失败: %v", err)
	}
	fmt.Printf("成功发送 %d 字节的数据\n", n)

	// 5. 读取服务器响应 (可选)
	// 设置读取超时时间，防止一直阻塞
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	n, err = conn.Read(buffer)
	if err != nil {
		if os.IsTimeout(err) {
			fmt.Println("读取响应超时（服务器可能没有返回数据）")
		} else if err == io.EOF {
			fmt.Println("服务器已关闭连接")
		} else {
			log.Fatalf("读取响应失败: %v", err)
		}
		return
	}

	// 打印响应内容（为了方便查看，将其转换为 Hex 字符串输出）
	fmt.Printf("收到服务器响应 (Hex): %X\n", buffer[:n])
}
