package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func canMake(t int64, a []int64, k int64) bool {
	n := int64(len(a))
	q := t / n
	r := int(t % n)
	need := int64(0)
	for i, v := range a {
		req := q
		if i < r {
			req++
		}
		if v < req {
			need += req - v
			if need > k {
				return false
			}
		}
	}
	return need <= k
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			sum += a[i]
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })

		low, high := int64(0), sum+k
		ans := int64(0)
		for low <= high {
			mid := (low + high) / 2
			if canMake(mid, a, k) {
				ans = mid
				low = mid + 1
			} else {
				high = mid - 1
			}
		}
		if ans >= int64(n) {
			fmt.Fprintln(writer, ans-int64(n)+1)
		} else {
			fmt.Fprintln(writer, 0)
		}
	}
}
