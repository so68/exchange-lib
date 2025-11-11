package main

import (
	"fmt"
	"time"
)

func main() {
	reloadFunc()
	select {}
}

func reloadFunc() {
	fmt.Println("=====>", time.Now().Format("2006-01-02 15:04:05"))
	time.AfterFunc(2*time.Second, reloadFunc)
}
