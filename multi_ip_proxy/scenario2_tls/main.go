package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"time"
)

// ---------------------------------------------------------
// 工程化核心思路：
// 1. 客户端拦截本地流量 -> 加上 TLS 加密外壳 -> 发送到海外服务端
// 2. 服务端监听 443 端口 -> 剥离 TLS 密码外壳 -> 将真实流量转发给目标网站 (这里用固定目标代替演示)
// ---------------------------------------------------------

const (
	// 假设我们在海外服务器的 4433 端口伪装 HTTPS 服务 (真实环境用 443)
	RemoteServerAddr = "127.0.0.1:4433"
	// 约定的鉴权头，用来区分是自己人，还是 GFW 的主动嗅探探测
	SecretToken = "Hello-Ziniao-Super-Secret"
)

// ========= 【服务端代码】运行在海外机房 =========
func startServer() {
	// 1. 加载 TLS 证书 (在生产中会使用 Let's Encrypt 等机构颁发的真实证书)
	// 这里为了演示，我们假设存在证书文件，或者为了让代码能直接运行，我们用极其简化的方法
	// *注意: 真实环境必须使用 tls.LoadX509KeyPair("server.crt", "server.key")
	cert, err := generateTempCert() // 动态生成临时证书以供演示不报错
	if err != nil {
		log.Fatalf("生成证书失败: %v", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", RemoteServerAddr, config)
	if err != nil {
		log.Fatalf("服务端启动失败: %v", err)
	}
	defer listener.Close()
	log.Println("[海外服务端] 已在", RemoteServerAddr, "开启 TLS 伪装监听，等待流量...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleServerConnection(conn)
	}
}

func handleServerConnection(clientConn net.Conn) {
	defer clientConn.Close()
	_ = clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// 2. 鉴权：读取请求头，看是不是自己人的客户端发来的
	buf := make([]byte, len(SecretToken))
	_, err := io.ReadFull(clientConn, buf)

	if err != nil || string(buf) != SecretToken {
		// 【防封锁核心】：主动探测防御 (Active Probing Defense)
		// 如果发现密码不对，绝不报错，直接伪装成一个正常的 Web 网站返回一个 HTML
		log.Println("[海外服务端] 发现不明探测请求，返回伪装的正常网页内容...")
		fakeHTML := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n<html><body><h1>Welcome to My Tech Blog</h1><p>Nothing here...</p></body></html>"
		_, _ = clientConn.Write([]byte(fakeHTML))
		return
	}

	// 3. 鉴权通过：恢复正常的超时时间
	_ = clientConn.SetReadDeadline(time.Time{})
	log.Println("[海外服务端] 客户端认证成功，正在代理流量至真实目标网站...")

	// 4. 作为代理，去请求真实的电商网站 (此处以 baidu 演示)
	targetConn, err := net.DialTimeout("tcp", "www.baidu.com:80", 5*time.Second)
	if err != nil {
		log.Println("连接真实目标失败:", err)
		return
	}
	defer targetConn.Close()

	// 5. 流量双向转发 (用 io.Copy 实现高性能的零拷贝转发)
	errChan := make(chan error, 2)
	go func() {
		_, err := io.Copy(targetConn, clientConn) // 将解密后的客户端请求，发送给目标网站
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(clientConn, targetConn) // 将目标网站返回的数据，加密通过 TLS 隧道传回给客户端
		errChan <- err
	}()

	<-errChan // 任意一方断开，结束该次代理
}

// ========= 【客户端代码】运行在紫鸟浏览器内部或用户本机 =========
func startClient() {
	log.Println("[本地客户端] 启动，监听本地 1080 端口接收浏览器流量...")
	listener, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		log.Fatalf("客户端启动失败: %v", err)
	}
	defer listener.Close()

	for {
		browserConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleBrowserConnection(browserConn)
	}
}

func handleBrowserConnection(browserConn net.Conn) {
	defer browserConn.Close()

	// 1. 客户端配置：由于服务端用的是临时生成的证书，这里跳过安全验证 (InsecureSkipVerify)
	// 【注意】生产环境中绝对不能设置为 true，且需要设置 ServerName 为伪装的域名防 SNI 审查
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "www.my-normal-website.com", // SNI 伪装
	}

	// 2. 拨号连接到海外服务端的 443(4433) 端口
	remoteConn, err := tls.Dial("tcp", RemoteServerAddr, tlsConfig)
	if err != nil {
		log.Println("[本地客户端] 连接海外服务器失败:", err)
		return
	}
	defer remoteConn.Close()

	// 3. 发送鉴权暗号，证明是自己人
	_, _ = remoteConn.Write([]byte(SecretToken))

	log.Println("[本地客户端] 隧道建立成功，正在将浏览器流量加密送往海外...")

	// 4. 将本地浏览器的明文数据注入 TLS 隧道加密；将海外回来的数据解密给浏览器
	errChan := make(chan error, 2)
	go func() {
		_, err := io.Copy(remoteConn, browserConn)
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(browserConn, remoteConn)
		errChan <- err
	}()

	<-errChan
}

// ========= 启动入口 =========
func main() {
	mode := flag.String("mode", "client", "运行模式: server 或 client")
	flag.Parse()

	if *mode == "server" {
		startServer()
	} else if *mode == "client" {
		startClient()
	} else {
		fmt.Println("请指定有效的运行模式。")
	}
}

// ============== 辅助函数：生成内存临时 TLS 证书 (仅供 Demo 直接运行用) ==============
func generateTempCert() (tls.Certificate, error) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return tls.X509KeyPair(certPEM, keyPEM)
}
