package main

import (
	"bufio"
	"fmt"
	"os"
)

func countAndSum(th int64, a []int64, k int64) (int64, int64) {
	var cnt, sum int64
	for _, ai := range a {
		if th < 1 {
			continue
		}
		t1 := ai
		if th < ai {
			t1 = th
		}
		cnt += t1
		sum += t1 * (t1 + 1) / 2

		upper := th - k
		if upper > ai {
			extra := upper - ai
			cnt += extra
			sum += (upper*(upper+1)/2 - ai*(ai+1)/2) + extra*k
		}
	}
	return cnt, sum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		var total int64
		var maxA int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			total += a[i]
			if a[i] > maxA {
				maxA = a[i]
			}
		}

		hi := maxA + k + total
		lo := int64(1)
		for lo < hi {
			mid := (lo + hi) >> 1
			cnt, _ := countAndSum(mid, a, k)
			if cnt >= total {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		th := lo

		cntLess, sumLess := countAndSum(th-1, a, k)
		_, sumTh := countAndSum(th, a, k)
		rem := total - cntLess
		ans := sumLess + rem*th

		// Safety: due to how we computed sumTh, ans should also equal sumTh - (cntTotal - total)*th
		_ = sumTh

		fmt.Fprintln(out, ans)
	}
}
