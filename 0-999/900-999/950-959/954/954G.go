package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var r int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &r, &k); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + a[i]
	}

	base := make([]int64, n+1)
	var mx int64
	for i := 1; i <= n; i++ {
		l := i - r
		if l < 1 {
			l = 1
		}
		rr := i + r
		if rr > n {
			rr = n
		}
		base[i] = prefix[rr] - prefix[l-1]
		if base[i] > mx {
			mx = base[i]
		}
	}

	low := int64(0)
	high := mx + k
	for low < high {
		mid := (low + high + 1) / 2
		if canReach(mid, base, n, r, k) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	fmt.Fprintln(writer, low)
}

func canReach(x int64, base []int64, n, r int, k int64) bool {
	diff := make([]int64, n+2)
	add := int64(0)
	left := k
	span := 2*r + 1

	for i := 1; i <= n; i++ {
		add += diff[i]
		cur := base[i] + add
		if cur < x {
			need := x - cur
			if need > left {
				return false
			}
			left -= need
			add += need
			end := i + span - 1
			if end <= n {
				diff[end+1] -= need
			}
		}
	}
	return true
}
