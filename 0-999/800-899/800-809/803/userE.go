package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	var s string
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	fmt.Fscan(in, &s)

	offset := k
	width := 2*k + 1

	parent := make([][]int, n+1)
	move := make([][]byte, n+1)
	for i := 0; i <= n; i++ {
		parent[i] = make([]int, width)
		move[i] = make([]byte, width)
		for j := 0; j < width; j++ {
			parent[i][j] = -2 // unreachable marker
		}
	}
	parent[0][offset] = -1 // start

	for i := 0; i < n; i++ {
		var cand []byte
		switch s[i] {
		case 'W', 'L', 'D':
			cand = []byte{s[i]}
		default:
			cand = []byte{'W', 'L', 'D'}
		}
		for d := 0; d < width; d++ {
			if parent[i][d] == -2 {
				continue
			}
			curDiff := d - offset
			for _, c := range cand {
				ndiff := curDiff
				switch c {
				case 'W':
					ndiff++
				case 'L':
					ndiff--
				case 'D':
					// no change
				}
				if ndiff < -k || ndiff > k {
					continue
				}
				if i != n-1 && abs(ndiff) == k {
					continue
				}
				nd := ndiff + offset
				if parent[i+1][nd] == -2 {
					parent[i+1][nd] = d
					move[i+1][nd] = c
				}
			}
		}
	}

	finalIdx := -1
	for d := 0; d < width; d++ {
		if parent[n][d] != -2 && abs(d-offset) == k {
			finalIdx = d
			break
		}
	}

	if finalIdx == -1 {
		fmt.Println("NO")
		return
	}

	res := make([]byte, n)
	idx := finalIdx
	for i := n; i > 0; i-- {
		res[i-1] = move[i][idx]
		idx = parent[i][idx]
	}
	fmt.Println(string(res))
}