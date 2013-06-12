package main

import (
	"fmt"
	"time"
)

func main() {
	var d = time.Duration(120e9)
	var i = int64(d)
	fmt.Println(d.String())
}
