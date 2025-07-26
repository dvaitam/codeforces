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
	sum := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		sum += x
	}
	fmt.Println(sum)
}
