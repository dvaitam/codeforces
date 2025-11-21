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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		a := make([]int, n)
		maxVal := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > maxVal {
				maxVal = a[i]
			}
		}

		m := n % k
		lenRow := n / k
		if m == 0 {
			m = k
		} else {
			lenRow++
		}

		rows := make([][]int, m)
		for r := 0; r < m; r++ {
			for idx := r; idx < n; idx += k {
				rows[r] = append(rows[r], a[idx])
			}
		}

		need := m/2 + 1
		lo, hi := 1, maxVal
		ans := 1
		for lo <= hi {
			mid := (lo + hi) / 2
			if can(rows, lenRow, need, mid) {
				ans = mid
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}

func can(rows [][]int, lenRow, need, threshold int) bool {
	m := len(rows)
	dpPrev := make([]int, lenRow)
	for j := 0; j < lenRow; j++ {
		if rows[0][j] >= threshold {
			dpPrev[j] = 1
		} else {
			dpPrev[j] = 0
		}
	}

	dpCurr := make([]int, lenRow)
	prefix := make([]int, lenRow)

	for r := 1; r < m; r++ {
		best := dpPrev[0]
		prefix[0] = best
		for j := 1; j < lenRow; j++ {
			if dpPrev[j] > best {
				best = dpPrev[j]
			}
			prefix[j] = best
		}
		for j := 0; j < lenRow; j++ {
			dpCurr[j] = prefix[j]
			if rows[r][j] >= threshold {
				dpCurr[j]++
			}
		}
		copy(dpPrev, dpCurr)
	}

	maxGood := dpPrev[0]
	for j := 1; j < lenRow; j++ {
		if dpPrev[j] > maxGood {
			maxGood = dpPrev[j]
		}
	}
	return maxGood >= need
}
