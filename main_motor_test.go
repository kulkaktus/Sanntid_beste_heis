package main

import (
	"./io"
	"fmt"
	"time"
)

func main() {
	for {
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Hello")
	}
}
