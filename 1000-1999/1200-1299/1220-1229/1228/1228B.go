package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

func modPow(a, e int) int {
	res := 1
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var h, w int
	if _, err := fmt.Fscan(in, &h, &w); err != nil {
		return
	}
	r := make([]int, h)
	for i := 0; i < h; i++ {
		fmt.Fscan(in, &r[i])
	}
	c := make([]int, w)
	for j := 0; j < w; j++ {
		fmt.Fscan(in, &c[j])
	}

	grid := make([][]int8, h)
	for i := 0; i < h; i++ {
		grid[i] = make([]int8, w)
	}

	// Apply row constraints
	for i := 0; i < h; i++ {
		for j := 0; j < r[i]; j++ {
			if grid[i][j] == -1 {
				fmt.Fprintln(out, 0)
				return
			}
			grid[i][j] = 1
		}
		if r[i] < w {
			if grid[i][r[i]] == 1 {
				fmt.Fprintln(out, 0)
				return
			}
			grid[i][r[i]] = -1
		}
	}

	// Apply column constraints
	for j := 0; j < w; j++ {
		for i := 0; i < c[j]; i++ {
			if grid[i][j] == -1 {
				fmt.Fprintln(out, 0)
				return
			}
			grid[i][j] = 1
		}
		if c[j] < h {
			if grid[c[j]][j] == 1 {
				fmt.Fprintln(out, 0)
				return
			}
			grid[c[j]][j] = -1
		}
	}

	freeCells := 0
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if grid[i][j] == 0 {
				freeCells++
			}
		}
	}

	ans := modPow(2, freeCells)
	fmt.Fprintln(out, ans)
}
