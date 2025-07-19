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
	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(in, &n)
		// As per solution: output four positive integers summing to n
		// 1, n-3, 1, 1
		fmt.Fprintf(out, "1 %d 1 1\n", n-3)
	}
}
