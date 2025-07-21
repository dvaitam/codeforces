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
		var l, r int64
		fmt.Fscan(in, &n, &l, &r)
		sum := int64(0)
		wins := 0
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			sum += v
			if sum > r {
				sum = 0
			} else if sum >= l {
				wins++
				sum = 0
			}
		}
		fmt.Fprintln(out, wins)
	}
}
