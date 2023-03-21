package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	var beginMem runtime.MemStats
	// 内存状态, Alloc: 已申请且仍在使用的字节
	runtime.ReadMemStats(&beginMem)
	fmt.Printf("KB: %v\n", beginMem.Alloc/1024)
	// 开始执行计时
	st := time.Now().UnixMilli()
	// go run code/user-submit/main.go
	cmd := exec.Command("go", "run", "code/user-submit/main.go")
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	// 输入测试用例
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(stdinPipe, "1 2\n")
	// 执行 go run code/user-submit/main.go 这条命令
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("stderr: ", stderr.String())
	}
	fmt.Println("Err: ", string(stderr.Bytes()))
	// 执行结果
	fmt.Println(out.String())
	fmt.Println("执行结果: ", out.String() == "3\n")
	// 内存和执行时间后处理
	var endMem runtime.MemStats
	runtime.ReadMemStats(&endMem)
	fmt.Printf("KB: %v\n", endMem.Alloc/1024)
	ed := time.Now().UnixMilli()
	fmt.Println("总耗时: ", ed-st)
}
