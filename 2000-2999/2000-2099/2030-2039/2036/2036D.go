package main

import (
	"bufio"
	"fmt"
	"os"
)

var pattern = []byte("1543")

func countPattern(seq []byte) int {
	if len(seq) < len(pattern) {
		return 0
	}
	cnt := 0
	ext := make([]byte, len(seq)+len(pattern)-1)
	copy(ext, seq)
	copy(ext[len(seq):], seq[:len(pattern)-1])
	for i := 0; i < len(seq); i++ {
		match := true
		for j := 0; j < len(pattern); j++ {
			if ext[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			cnt++
		}
	}
	return cnt
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var line string
			fmt.Fscan(in, &line)
			grid[i] = []byte(line)
		}

		ans := 0
		top, left := 0, 0
		bottom, right := n-1, m-1
		for top <= bottom && left <= right {
			if top == bottom || left == right {
				break
			}
			layerLen := 2 * (bottom - top + right - left)
			seq := make([]byte, 0, layerLen)

			for j := left; j <= right; j++ {
				seq = append(seq, grid[top][j])
			}
			for i := top + 1; i <= bottom-1; i++ {
				seq = append(seq, grid[i][right])
			}
			if bottom > top {
				for j := right; j >= left; j-- {
					seq = append(seq, grid[bottom][j])
				}
			}
			if right > left {
				for i := bottom - 1; i >= top+1; i-- {
					seq = append(seq, grid[i][left])
				}
			}

			ans += countPattern(seq)

			top++
			left++
			bottom--
			right--
		}
		fmt.Fprintln(out, ans)
	}
}
