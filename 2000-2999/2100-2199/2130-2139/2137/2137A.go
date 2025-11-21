package main

import (
	"bufio"
	"fmt"
	"os"
)

func reverseCollatz(x int) int {
	return 2 * x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var k, x int
		fmt.Fscan(in, &k, &x)
		result := x
		for i := 0; i < k; i++ {
			result = reverseCollatz(result)
		}
		fmt.Fprintln(out, result)
	}
}
