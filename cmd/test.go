package main

import (
	"fmt"
	"time"
)

func main() {
	a := time.Now()
	b := time.Now()

	fmt.Println(a.Sub(b))

	fmt.Println(time.Now())
	time.Sleep(a.Sub(b))
	fmt.Println(time.Now())
}


