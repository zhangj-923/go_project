package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// IPProxyManager 管理不同本地出网 IP 的 HttpClient
// 工程化要点 1：复用 http.Client 和 Transport，避免每次请求都创建，浪费资源且无法利用连接池 (Keep-Alive)
type IPProxyManager struct {
	clients map[string]*http.Client
	mu      sync.RWMutex
}

func NewIPProxyManager() *IPProxyManager {
	return &IPProxyManager{
		clients: make(map[string]*http.Client),
	}
}

// GetClient 根据传入的本地出网 IP 获取一个配置好的 HttpClient
func (m *IPProxyManager) GetClient(localIPStr string) (*http.Client, error) {
	// 先尝试读锁，看是否已经缓存了该 IP 的 Client
	m.mu.RLock()
	client, exists := m.clients[localIPStr]
	m.mu.RUnlock()
	if exists {
		return client, nil
	}

	// 验证 IP 格式是否合法
	localIP := net.ParseIP(localIPStr)
	if localIP == nil {
		return nil, fmt.Errorf("无效的本地 IP 地址: %s", localIPStr)
	}

	// 加写锁创建新的 Client
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查锁定 (Double-Check Locking)，防止并发时多次创建
	if client, exists = m.clients[localIPStr]; exists {
		return client, nil
	}

	// 工程化要点 2：定制 net.Dialer 强制绑定 LocalAddr
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   localIP,
			Port: 0, // 0 表示由操作系统随机分配一个可用的临时端口
		},
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second, // 开启 TCP Keep-Alive
	}

	// 工程化要点 3：定制 http.Transport
	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,              // 整个客户端的最大空闲连接数
		IdleConnTimeout:       90 * time.Second, // 空闲连接超时时间
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   10, // 每个目标主机的最大空闲连接数，控制连接池大小
	}

	newClient := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second, // 整个 HTTP 请求的超时时间
	}

	// 缓存起来
	m.clients[localIPStr] = newClient
	return newClient, nil
}

// FetchData 使用指定的本地 IP 抓取目标网址的数据
func (m *IPProxyManager) FetchData(ctx context.Context, localIP string, targetURL string) {
	client, err := m.GetClient(localIP)
	if err != nil {
		log.Printf("[Error] 获取 %s 的 Client 失败: %v\n", localIP, err)
		return
	}

	// 使用带有 Context 的 Request，方便进行超时控制或主动取消
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		log.Printf("[Error] 创建请求失败: %v\n", err)
		return
	}

	// 伪装一些基础的 Header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Error] IP %s 请求 %s 失败: %v\n", localIP, targetURL, err)
		return
	}
	defer resp.Body.Close()

	// 读取前 100 个字节用于演示
	bodyBytes := make([]byte, 100)
	n, _ := io.ReadFull(resp.Body, bodyBytes)

	log.Printf("[Success] 本地出口 IP: %s | 目标: %s | 状态码: %d | 耗时: %v | 响应片段: %s...\n",
		localIP, targetURL, resp.StatusCode, time.Since(start), string(bodyBytes[:n]))
}

func main() {
	manager := NewIPProxyManager()

	// 测试目标地址，如果目标是亚马逊，可以使用类似 https://www.amazon.com
	targetURL := "https://httpbin.org/ip" // 这个 API 会返回你请求它的真实出口 IP

	// 模拟多个并发的店铺环境请求
	// 【注意】如果要真实运行成功，你的电脑或者服务器上必须真的配置了下面这些 IP
	// 否则底层的 bind() 系统调用会报 "bind: can't assign requested address" 错误
	mockShopIPs := []string{
		"192.168.1.100", // 假设这是分配给店铺 A 的出网 IP
		"192.168.1.101", // 假设这是分配给店铺 B 的出网 IP
		"127.0.0.1",     // 本机回环地址测试（一定会成功，但不走外网）
	}

	var wg sync.WaitGroup
	ctx := context.Background()

	for _, ip := range mockShopIPs {
		wg.Add(1)
		go func(localIP string) {
			defer wg.Done()
			manager.FetchData(ctx, localIP, targetURL)
		}(ip)
	}

	wg.Wait()
	fmt.Println("所有模拟请求执行完毕。")
}
