package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	fib := make([]int64, 21)
	fib[1] = 1
	for i := 2; i < len(fib); i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		fmt.Fprintln(out, fib[n])
	}
}
