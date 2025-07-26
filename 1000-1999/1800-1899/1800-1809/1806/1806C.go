package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int64, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Fscan(in, &p[i])
		}
		if n == 1 {
			ans := absInt64(p[0] - p[1])
			fmt.Fprintln(out, ans)
			continue
		}
		if n == 2 {
			var sum0, sum2, sumMinus int64
			bestDelta := int64(1 << 62)
			for i := 0; i < 4; i++ {
				x := p[i]
				sum0 += absInt64(x)
				sum2 += absInt64(x - 2)
				sumMinus += absInt64(x + 1)
				delta := absInt64(x-2) - absInt64(x+1)
				if delta < bestDelta {
					bestDelta = delta
				}
			}
			cand3 := sumMinus + bestDelta
			ans := sum0
			if sum2 < ans {
				ans = sum2
			}
			if cand3 < ans {
				ans = cand3
			}
			fmt.Fprintln(out, ans)
			continue
		}
		var sum int64
		for _, v := range p {
			sum += absInt64(v)
		}
		fmt.Fprintln(out, sum)
	}
}
