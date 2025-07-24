package main

import (
	"bufio"
	"fmt"
	"os"
)

type Box struct {
	minR, maxR int
	minC, maxC int
	r0         int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &mat[i][j])
		}
	}

	N := n * n
	const inf = int(1e9)
	minR := make([]int, N+1)
	maxR := make([]int, N+1)
	minC := make([]int, N+1)
	maxC := make([]int, N+1)
	present := make([]bool, N+1)
	for i := 0; i <= N; i++ {
		minR[i] = inf
		minC[i] = inf
		maxR[i] = -1
		maxC[i] = -1
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := mat[i][j]
			if !present[v] {
				present[v] = true
				minR[v] = i
				maxR[v] = i
				minC[v] = j
				maxC[v] = j
			} else {
				if i < minR[v] {
					minR[v] = i
				}
				if i > maxR[v] {
					maxR[v] = i
				}
				if j < minC[v] {
					minC[v] = j
				}
				if j > maxC[v] {
					maxC[v] = j
				}
			}
		}
	}

	boxes := make([]Box, 0)
	for v := 1; v <= N; v++ {
		if present[v] {
			r0 := max(maxR[v]-minR[v]+1, maxC[v]-minC[v]+1)
			boxes = append(boxes, Box{minR[v], maxR[v], minC[v], maxC[v], r0})
		}
	}
	m := len(boxes)
	if m == k {
		fmt.Fprintln(out, 0)
		return
	}
	if m < k {
		fmt.Fprintln(out, k-m)
		return
	}

	need1 := m - k
	need2 := m - k + 1

	for r := 1; r <= n; r++ {
		size := n - r + 1
		diff := make([][]int, size+1)
		for i := range diff {
			diff[i] = make([]int, size+1)
		}
		for _, b := range boxes {
			if b.r0 > r {
				continue
			}
			sr := b.maxR - r + 1
			er := b.minR
			sc := b.maxC - r + 1
			ec := b.minC
			if sr < 0 {
				sr = 0
			}
			if sc < 0 {
				sc = 0
			}
			if er > n-r {
				er = n - r
			}
			if ec > n-r {
				ec = n - r
			}
			if sr <= er && sc <= ec {
				diff[sr][sc]++
				diff[er+1][sc]--
				diff[sr][ec+1]--
				diff[er+1][ec+1]++
			}
		}
		for i := 0; i <= size; i++ {
			for j := 1; j <= size; j++ {
				diff[i][j] += diff[i][j-1]
			}
		}
		for i := 1; i <= size; i++ {
			for j := 0; j <= size; j++ {
				diff[i][j] += diff[i-1][j]
			}
		}
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				c := diff[i][j]
				if c == need1 || c == need2 {
					fmt.Fprintln(out, 1)
					return
				}
			}
		}
	}

	fmt.Fprintln(out, 2)
}
