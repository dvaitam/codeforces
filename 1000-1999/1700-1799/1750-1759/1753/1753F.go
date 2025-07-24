package main

import (
	"bufio"
	"fmt"
	"os"
)

type cell struct {
	weights []int
	set     map[int]struct{}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func addCell(c *cell, posCount, negCount []int, posPresent, negPresent []bool, rFirst, nFirst *int, limit int) {
	for _, w := range c.weights {
		if w > 0 {
			if w > limit {
				continue
			}
			if posCount[w] == 0 {
				posPresent[w] = true
			}
			posCount[w]++
			if w == *rFirst {
				for *rFirst <= limit && posPresent[*rFirst] {
					*rFirst++
				}
			}
		} else {
			idx := -w
			if idx > limit {
				continue
			}
			if negCount[idx] == 0 {
				negPresent[idx] = true
			}
			negCount[idx]++
			if idx == *nFirst {
				for *nFirst <= limit && negPresent[*nFirst] {
					*nFirst++
				}
			}
		}
	}
}

func removeCell(c *cell, posCount, negCount []int, posPresent, negPresent []bool, rFirst, nFirst *int) {
	for _, w := range c.weights {
		if w > 0 {
			if posCount[w] > 0 {
				posCount[w]--
				if posCount[w] == 0 {
					posPresent[w] = false
					if w < *rFirst {
						*rFirst = w
					}
				}
			}
		} else {
			idx := -w
			if negCount[idx] > 0 {
				negCount[idx]--
				if negCount[idx] == 0 {
					negPresent[idx] = false
					if idx < *nFirst {
						*nFirst = idx
					}
				}
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k, t int
	if _, err := fmt.Fscan(reader, &n, &m, &k, &t); err != nil {
		return
	}

	trans := false
	if n > m {
		n, m = m, n
		trans = true
	}

	limit := t - 1
	cells := make([][]cell, n)
	for i := range cells {
		cells[i] = make([]cell, m)
	}

	for i := 0; i < k; i++ {
		var x, y, w int
		fmt.Fscan(reader, &x, &y, &w)
		if trans {
			x, y = y, x
		}
		if abs(w) > limit {
			continue
		}
		c := &cells[x-1][y-1]
		if c.set == nil {
			c.set = make(map[int]struct{})
		}
		if _, ok := c.set[w]; !ok {
			c.set[w] = struct{}{}
			c.weights = append(c.weights, w)
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			cells[i][j].set = nil
		}
	}

	if t <= 1 {
		ans := 0
		minnm := n
		if m < n {
			minnm = m
		}
		for L := 1; L <= minnm; L++ {
			ans += (n - L + 1) * (m - L + 1)
		}
		fmt.Fprintln(writer, ans)
		return
	}

	posCount := make([]int, limit+1)
	negCount := make([]int, limit+1)
	posPresent := make([]bool, limit+1)
	negPresent := make([]bool, limit+1)

	minnm := n
	if m < n {
		minnm = m
	}
	ans := 0
	for L := 1; L <= minnm; L++ {
		for top := 0; top <= n-L; top++ {
			for i := 1; i <= limit; i++ {
				posCount[i] = 0
				negCount[i] = 0
				posPresent[i] = false
				negPresent[i] = false
			}
			rFirst, nFirst := 1, 1
			for i := 0; i < L; i++ {
				for j := 0; j < L; j++ {
					addCell(&cells[top+i][j], posCount, negCount, posPresent, negPresent, &rFirst, &nFirst, limit)
				}
			}
			teamSize := rFirst + nFirst - 1
			if teamSize >= t {
				ans++
			}
			for left := 1; left <= m-L; left++ {
				for i := 0; i < L; i++ {
					removeCell(&cells[top+i][left-1], posCount, negCount, posPresent, negPresent, &rFirst, &nFirst)
					addCell(&cells[top+i][left+L-1], posCount, negCount, posPresent, negPresent, &rFirst, &nFirst, limit)
				}
				if rFirst <= limit {
					for rFirst <= limit && posPresent[rFirst] {
						rFirst++
					}
				}
				if nFirst <= limit {
					for nFirst <= limit && negPresent[nFirst] {
						nFirst++
					}
				}
				teamSize = rFirst + nFirst - 1
				if teamSize >= t {
					ans++
				}
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
