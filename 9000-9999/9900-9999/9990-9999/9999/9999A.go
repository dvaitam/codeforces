package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	maxVal := 0
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] > maxVal {
			maxVal = a[i]
		}
	}
	// determine bit width
	B := 0
	for (1 << B) <= maxVal {
		B++
	}
	if B == 0 {
		B = 1
	}

	last := make([]int, B)
	L := make([]int, n+1)
	for i := 1; i <= n; i++ {
		maxIdx := 0
		for b := 0; b < B; b++ {
			if (a[i]>>b)&1 == 0 {
				if last[b] > maxIdx {
					maxIdx = last[b]
				}
			}
		}
		L[i] = maxIdx + 1
		for b := 0; b < B; b++ {
			if (a[i]>>b)&1 == 1 {
				last[b] = i
			}
		}
	}

	next := make([]int, B)
	R := make([]int, n+1)
	for b := 0; b < B; b++ {
		next[b] = n + 1
	}
	for i := n; i >= 1; i-- {
		minIdx := n + 1
		for b := 0; b < B; b++ {
			if (a[i]>>b)&1 == 0 {
				if next[b] < minIdx {
					minIdx = next[b]
				}
			}
		}
		R[i] = minIdx - 1
		for b := 0; b < B; b++ {
			if (a[i]>>b)&1 == 1 {
				next[b] = i
			}
		}
	}

	prevSame := make(map[int]int)
	var ans int64
	for i := 1; i <= n; i++ {
		start := L[i]
		if p, ok := prevSame[a[i]]; ok && p+1 > start {
			start = p + 1
		}
		left := i - start + 1
		right := R[i] - i + 1
		if left > 0 && right > 0 {
			ans += int64(left) * int64(right)
		}
		prevSame[a[i]] = i
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
