package main

import (
	"fmt"
	"time"
)

// 两数之和
func main() {
	var a, b int
	fmt.Scanln(&a, &b)
	fmt.Println(a + b + 1)
	time.Sleep(time.Second * 10)
}
