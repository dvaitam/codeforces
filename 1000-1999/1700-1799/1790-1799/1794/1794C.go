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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		prefix := make([]float64, n+1)
		for i := 1; i <= n; i++ {
			prefix[i] = prefix[i-1] + math.Log(float64(a[i-1]))
		}
		logFact := make([]float64, n+1)
		for i := 1; i <= n; i++ {
			logFact[i] = logFact[i-1] + math.Log(float64(i))
		}

		ans := make([]int, n)
		d := 0
		for k := 1; k <= n; k++ {
			for d < k {
				scoreD := prefix[k] - prefix[k-d] - logFact[d]
				scoreD1 := prefix[k] - prefix[k-(d+1)] - logFact[d+1]
				if scoreD1 >= scoreD-1e-12 {
					d++
				} else {
					break
				}
			}
			ans[k-1] = d
		}

		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
