package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	var k int64
	var x, y int
	fmt.Fscan(in, &n, &m, &k, &x, &y)

	counts := make([][]int64, n)
	for i := range counts {
		counts[i] = make([]int64, m)
	}

	if n == 1 {
		full := k / int64(m)
		rem := k % int64(m)
		for j := 0; j < m; j++ {
			counts[0][j] = full
		}
		for j := 0; j < int(rem); j++ {
			counts[0][j]++
		}
	} else {
		cycleRows := int64(2*n - 2)
		cycleQuestions := cycleRows * int64(m)
		fullCycles := k / cycleQuestions
		// base counts
		for i := 0; i < n; i++ {
			var visits int64
			if i == 0 || i == n-1 {
				visits = fullCycles
			} else {
				visits = fullCycles * 2
			}
			for j := 0; j < m; j++ {
				counts[i][j] = visits
			}
		}
		rem := k % cycleQuestions
		if rem > 0 {
			// simulate remaining questions
			// build row order
			rows := make([]int, 0, 2*n-2)
			for i := 0; i < n; i++ {
				rows = append(rows, i)
			}
			for i := n - 2; i >= 1; i-- {
				rows = append(rows, i)
			}
			for _, r := range rows {
				if rem == 0 {
					break
				}
				if rem >= int64(m) {
					for j := 0; j < m; j++ {
						counts[r][j]++
					}
					rem -= int64(m)
				} else {
					for j := 0; j < int(rem); j++ {
						counts[r][j]++
					}
					rem = 0
				}
			}
		}
	}

	var maxv int64
	minv := counts[0][0]
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if counts[i][j] > maxv {
				maxv = counts[i][j]
			}
			if counts[i][j] < minv {
				minv = counts[i][j]
			}
		}
	}
	fmt.Printf("%d %d %d\n", maxv, minv, counts[x-1][y-1])
}
