package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func fib(n int) int64 {
	if n <= 1 {
		return int64(n)
	}
	a, b := int64(0), int64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanWords)
	in.Scan()
	t, _ := strconv.Atoi(in.Text())
	for i := 0; i < t; i++ {
		in.Scan()
		n, _ := strconv.Atoi(in.Text())
		fmt.Println(fib(n))
	}
}
