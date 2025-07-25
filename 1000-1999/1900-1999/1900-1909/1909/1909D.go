package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}

		g := int64(0)
		allZero := true
		b := make([]int64, n)
		for i, v := range a {
			b[i] = v - k
			if b[i] != 0 {
				allZero = false
			}
			if b[i] < 0 {
				b[i] = -b[i]
			}
			g = gcd(g, b[i])
		}
		if allZero {
			fmt.Fprintln(out, 0)
			continue
		}

		// collect divisors of g
		divisors := make([]int64, 0)
		for d := int64(1); d*d <= g; d++ {
			if g%d == 0 {
				divisors = append(divisors, d)
				if d != g/d {
					divisors = append(divisors, g/d)
				}
			}
		}

		ans := int64(-1)
		// try both positive and negative divisors
		for _, d := range divisors {
			for _, div := range []int64{d, -d} {
				Tval := k + div
				if Tval <= 0 {
					continue
				}
				ok := true
				sumQ := int64(0)
				for i := 0; i < n; i++ {
					if (a[i]-k)%div != 0 {
						ok = false
						break
					}
					q := (a[i] - k) / div
					if q <= 0 {
						ok = false
						break
					}
					sumQ += q
				}
				if ok {
					ops := sumQ - int64(n)
					if ans == -1 || ops < ans {
						ans = ops
					}
				}
			}
		}

		fmt.Fprintln(out, ans)
	}
}
