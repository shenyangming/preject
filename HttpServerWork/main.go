package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/golang/glog"
)

func main() {
	mux := http.NewServeMux()
	//1.接收客户端ClientHeader获取header
	mux.HandleFunc("/header", ResponseHeader)

	//2.获取系统环境变量VERSON配置
	mux.HandleFunc("/verson", ResponseVersion)

	//3.接收客户端Clientip获取IP,返回码
	mux.HandleFunc("/IP", remoteAddr)

	//4.访问localhost/healthz返回200
	mux.HandleFunc("/healthz", healthz)

	// 设置服务器
	server := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}
	// 设置服务器监听请求端口
	server.ListenAndServe()
	err := server.ListenAndServe().Error()
	if err != "" {
		log.Fatal(err)
	}

}

/*1.接收客户端 request，并将 request 中带的 header 写入 response header*/
func ResponseHeader(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--服务器开启,开始写入request header--")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
		w.Header().Add(k, v[0])
	}
	io.WriteString(w, "===================Details of the http Response header:============\n")
	for k, v := range w.Header() {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}

/*2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header*/
func ResponseVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--服务器开启,开始写入读取环境变量--")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================下面是系统环境变量VERSION的配置============\n")

	err := os.Setenv("VERSION", "1.6.0") //临时设置 系统环境变量
	if err != nil {
		fmt.Println(err.Error())
	}

	envs := os.Getenv("VERSION")
	w.Header().Add("VERSION", envs) //写入response header

	/*遍历response*/
	for k, v := range w.Header() {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

	/*envs := os.Environ()
	遍历环境变量并写入请求-
	for _, v := range envs {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			continue
		} else {
			w.Header().Add(parts[0], parts[1])
		}
	}
	/*遍历变量并写入请求-
	for k, v := range w.Header() {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}*/
}

/*3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出*/
func remoteAddr(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--服务器开启,开始读取IP和返回码--")

	/*定义标签*/
	const (
		XForwardedFor = "X-Forwarded-For"
		XRealIP       = "X-Real-IP"
	)
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "192.168.0.149"
	}
	fmt.Println(http.StatusOK)

	/*存入日志*/
	flag.Parse()
	defer glog.Flush()
	flag.Set("v", "4")
	glog.V(2).Info("客户端的IP: ", remoteAddr, "HTTP返回码：", http.StatusOK)

	/*标准输出*/
	fmt.Printf("客户端的IP:%s,HTTP返回码：%v", remoteAddr, http.StatusOK)

	/*服务器端输出*/
	io.WriteString(w, fmt.Sprintf("客户端的IP:%s,HTTP返回码：%v", remoteAddr, http.StatusOK))
}

/*4. 当访问 localhost/healthz 时，应返回200*/
func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--服务器开启,开始返回200--")
	io.WriteString(w, "200\n")
}
