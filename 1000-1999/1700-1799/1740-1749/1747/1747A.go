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
	fmt.Fscan(in, &t)
	for t > 0 {
		t--
		var n int
		fmt.Fscan(in, &n)
		var sum int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			sum += x
		}
		if sum < 0 {
			sum = -sum
		}
		fmt.Fprintln(out, sum)
	}
}
