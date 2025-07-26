package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b int64
		fmt.Fscan(in, &a, &b)
		if a > b {
			a, b = b, a
		}
		ans := a + b
		limit := int64(math.Sqrt(float64(b))) + 2
		for i := int64(1); i <= limit; i++ {
			cands := []int64{i, a / i, a/i + 1, b / i, b/i + 1}
			for _, s := range cands {
				if s <= 0 {
					continue
				}
				val := (s - 1) + (a+s-1)/s + (b+s-1)/s
				if val < ans {
					ans = val
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
