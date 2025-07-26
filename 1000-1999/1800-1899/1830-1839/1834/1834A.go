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
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		neg := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x == -1 {
				neg++
			}
		}
		sum := n - 2*neg
		ops := 0
		if neg%2 == 1 {
			ops++
			neg--
			sum += 2
		}
		if sum < 0 {
			deficit := -sum
			pairs := (deficit + 3) / 4
			ops += pairs * 2
		}
		fmt.Fprintln(out, ops)
	}
}
