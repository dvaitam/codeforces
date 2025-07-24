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

	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		sum += a[i]
	}

	limit := int64(math.Sqrt(float64(sum+k))) + 2
	ans := int64(0)

	check := func(d int64) bool {
		var waste int64
		for _, v := range a {
			waste += ((v+d-1)/d)*d - v
			if waste > k {
				return false
			}
		}
		return waste <= k
	}

	for d := int64(1); d <= limit; d++ {
		if check(d) {
			if d > ans {
				ans = d
			}
		}
	}
	for q := int64(1); q <= limit; q++ {
		d := (sum + k) / q
		if d > ans && check(d) {
			ans = d
		}
	}

	fmt.Fprintln(out, ans)
}
