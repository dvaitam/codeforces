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

	var n int
	fmt.Fscan(in, &n)

	var sum int64
	diff := 0
	for i := 1; i <= n; i++ {
		var a int64
		fmt.Fscan(in, &a)
		sum += a
		if a&1 == 1 {
			if i&1 == 1 {
				diff++
			} else {
				diff--
			}
		}
	}
	if diff < 0 {
		diff = -diff
	}
	ans := (sum - int64(diff)) / 2
	fmt.Fprintln(out, ans)
}
