package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	const INF int = int(1e18)

	for ; T > 0; T-- {
		var n, d int
		fmt.Fscan(reader, &n, &d)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, n+2)
		b[0] = 0
		for i := 1; i <= n; i++ {
			b[i] = a[i-1]
		}
		b[n+1] = d + 1
		g := make([]int, n+1)
		for i := 1; i <= n+1; i++ {
			g[i-1] = b[i] - b[i-1] - 1
		}
		m := n
		prefMin := make([]int, m+1)
		prefMax := make([]int, m+1)
		prefMin[0] = INF
		for i := 0; i < m; i++ {
			val := g[i]
			if val < prefMin[i] {
				prefMin[i+1] = val
			} else {
				prefMin[i+1] = prefMin[i]
			}
			if prefMax[i] > val {
				prefMax[i+1] = prefMax[i]
			} else {
				prefMax[i+1] = val
			}
		}
		suffMin := make([]int, m+2)
		suffMax := make([]int, m+2)
		suffMin[m] = INF
		for i := m - 1; i >= 0; i-- {
			val := g[i]
			if val < suffMin[i+1] {
				suffMin[i] = val
			} else {
				suffMin[i] = suffMin[i+1]
			}
			if suffMax[i+1] > val {
				suffMax[i] = suffMax[i+1]
			} else {
				suffMax[i] = val
			}
		}
		baseline := suffMin[0]
		ans := baseline
		for i := 0; i < n; i++ {
			newGap := g[i] + g[i+1] + 1
			leftMin := prefMin[i]
			rightIdx := i + 2
			rightMin := INF
			if rightIdx <= m {
				rightMin = suffMin[rightIdx]
			}
			minWithout := min(leftMin, rightMin)
			minWithout = min(minWithout, newGap)

			leftMax := prefMax[i]
			rightMax := 0
			if rightIdx <= m {
				rightMax = suffMax[rightIdx]
			}
			maxWithout := max(leftMax, rightMax)
			maxWithout = max(maxWithout, newGap)

			lastExam := a[n-1]
			if i == n-1 {
				if n >= 2 {
					lastExam = a[n-2]
				} else {
					lastExam = 0
				}
			}
			trailing := d - lastExam - 1
			if trailing < 0 {
				trailing = 0
			}

			cand1 := min(minWithout, (maxWithout-1)/2)
			if cand1 > ans {
				ans = cand1
			}
			cand2 := min(minWithout, trailing)
			if cand2 > ans {
				ans = cand2
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
