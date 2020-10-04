package main

import (
"flag"
"fmt"
"runtime"
"strings"
"go-stress-testing/model"
"go-stress-testing/server"
)

type array []string

func (a *array) String() string {
	return fmt.Sprint(*a)
}

func (a *array) Set(s string) error {
	*a = append(*a, s)

	return nil
}

func main() {

	runtime.GOMAXPROCS(1)

	var (
		concurrency uint64
		totalNumber uint64
		debugStr    string
		requestUrl  string
		path        string
		verify      string
		headers     array
		body        string
	)

	flag.Uint64Var(&concurrency, "c", 1, "并发数")
	flag.Uint64Var(&totalNumber, "n", 1, "请求数(单个并发/协程)")
	flag.StringVar(&debugStr, "d", "false", "调试模式")
	flag.StringVar(&requestUrl, "u", "", "压测地址")
	flag.StringVar(&path, "p", "", "curl文件路径")
	flag.StringVar(&verify, "v", "", "验证方法 htt")
	flag.Var(&headers, "H", "自定义头信息传递给服务器")
	flag.StringVar(&body, "data", "", "HTTP POST方式传送数据")

	flag.Parse()
	if concurrency == 0 || totalNumber == 0 || (requestUrl == "" && path == "") {
		fmt.Printf("压测地址 \n")
		fmt.Printf("当前请求参数: -c %d -n %d -d %v -u %s \n", concurrency, totalNumber, debugStr, requestUrl)

		flag.Usage()

		return
	}

	debug := strings.ToLower(debugStr) == "true"
	request, err := model.NewRequest(requestUrl, verify, 0, debug, path, headers, body)
	if err != nil {
		fmt.Printf("参数不合法 %v \n", err)

		return
	}

	fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", concurrency, totalNumber)
	request.Print()

	server.Dispose(concurrency, totalNumber, request)

	return
}
