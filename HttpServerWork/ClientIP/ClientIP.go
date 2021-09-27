package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("启动客户端")
	// 使用Get方法获取服务器响应包数据
	resp, err := http.Get("http://192.168.0.149:80/IP")
	if err != nil {
		fmt.Println("Get err:", err)
		return
	}
	defer resp.Body.Close()

	// 获取服务器端读到的数据
	fmt.Println("Status = ", resp.Status)         // 状态
	fmt.Println("StatusCode = ", resp.StatusCode) // 状态码
	fmt.Println("Header = ", resp.Header)         // 响应头部
	fmt.Println("Body = ", resp.Body)             // 响应包体

	//读取body内的内容
	content, err := ioutil.ReadAll(resp.Body)

	// 打印从body中读到的所有内容
	fmt.Println("result = ", string(content))

}
