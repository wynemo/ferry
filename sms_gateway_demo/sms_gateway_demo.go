package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	queryParams := r.URL.Query()

	// 打印参数到服务器日志
	fmt.Println("收到的 GET 请求参数:")
	for key, values := range queryParams {
		fmt.Printf("%s: %s\n", key, values)
	}

	// 返回响应
	fmt.Fprintf(w, "收到的参数:\n")
	for key, values := range queryParams {
		fmt.Fprintf(w, "%s: %s\n", key, values)
	}
}

func main() {
	// 注册路由
	http.HandleFunc("/", handler)

	// 启动服务器
	fmt.Println("服务器启动，监听端口 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}
