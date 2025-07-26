package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	res := 1
	for i := 2; i <= n; i++ {
		res *= i
	}
	fmt.Println(res)
}
