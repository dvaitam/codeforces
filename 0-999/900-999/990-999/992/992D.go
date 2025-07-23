package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	onesCount := 0
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		prefix[i+1] = prefix[i] + int64(a[i])
		if a[i] == 1 {
			onesCount++
		}
	}

	positions := make([]int, 0)
	for i := 0; i < n; i++ {
		if a[i] > 1 {
			positions = append(positions, i)
		}
	}

	var ans int64
	if k == 1 {
		ans += int64(onesCount)
	}

	limit := int64(2000000000000000000) // 2e18
	m := len(positions)
	for i := 0; i < m; i++ {
		prod := int64(1)
		for j := i; j < m && j < i+60; j++ {
			val := int64(a[positions[j]])
			if prod > limit/val {
				break
			}
			prod *= val
			L := positions[i]
			R := positions[j]

			var leftOnes, rightOnes int
			if i == 0 {
				leftOnes = L
			} else {
				leftOnes = L - positions[i-1] - 1
			}
			if j == m-1 {
				rightOnes = n - R - 1
			} else {
				rightOnes = positions[j+1] - R - 1
			}

			sumSeg := prefix[R+1] - prefix[L]

			if prod%int64(k) != 0 {
				continue
			}
			target := prod / int64(k)
			need := target - sumSeg
			if need < 0 || need > int64(leftOnes+rightOnes) {
				continue
			}
			low := int64(0)
			if need-int64(rightOnes) > 0 {
				low = need - int64(rightOnes)
			}
			high := need
			if int64(leftOnes) < high {
				high = int64(leftOnes)
			}
			if low <= high {
				ans += high - low + 1
			}
		}
	}

	fmt.Fprintln(out, ans)
}
