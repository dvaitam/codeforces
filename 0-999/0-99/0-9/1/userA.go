package main

import (
	"fmt"
)

func main() {
	var n, m, a int64
	fmt.Scan(&n, &m, &a)
	
	rows := (n + a - 1) / a
	cols := (m + a - 1) / a
	
	fmt.Println(rows * cols)
}