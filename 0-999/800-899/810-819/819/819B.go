package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}

	diff := make([]int64, n+2)
	add := func(l, r int, val int64) {
		if l < 1 {
			l = 1
		}
		if r > n {
			r = n
		}
		if l > r {
			return
		}
		diff[l] += val
		if r+1 <= n {
			diff[r+1] -= val
		}
	}

	for i := 1; i <= n; i++ {
		v := p[i-1]
		w := n - i + 1

		start := 1
		end := n - i
		if start <= end {
			if v <= i {
				add(start, end, 1)
			} else if v <= n-1 {
				add(start, start+v-i-1, -1)
				add(start+v-i, end, 1)
			} else {
				add(start, end, -1)
			}
		}

		add(w, w, int64(2*v-n-1))

		start = w + 1
		end = w + i - 1
		if start <= end {
			if v > i-1 {
				add(start, end, -1)
			} else {
				add(start, start+v-2, -1)
				add(start+v-1, end, 1)
			}
		}
	}

	for k := 1; k <= n; k++ {
		diff[k] += diff[k-1]
	}

	cur := int64(0)
	for i := 0; i < n; i++ {
		d := int64(p[i] - (i + 1))
		if d < 0 {
			d = -d
		}
		cur += d
	}

	bestVal := cur
	bestIdx := 0
	for k := 1; k < n; k++ {
		cur += diff[k]
		if cur < bestVal {
			bestVal = cur
			bestIdx = k
		}
	}

	fmt.Println(bestVal, bestIdx)
}
