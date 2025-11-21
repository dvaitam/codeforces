package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const mod int64 = 998244353

func isqrt(x int64) int64 {
	if x <= 0 {
		return 0
	}
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func rangeSum(prefix []int64, l, r int) int64 {
	if l < 0 {
		l = 0
	}
	if r > len(prefix)-1 {
		r = len(prefix) - 1
	}
	if l >= r {
		return 0
	}
	val := prefix[r] - prefix[l]
	val %= mod
	if val < 0 {
		val += mod
	}
	return val
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, d int
		fmt.Fscan(in, &n, &m, &d)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}

		dpIn := make([]int64, m)
		dpLeave := make([]int64, m)
		prefix := make([]int64, m+1)

		bottom := grid[n-1]
		hasBottom := false
		for c := 0; c < m; c++ {
			if bottom[c] == 'X' {
				dpIn[c] = 1
				hasBottom = true
			} else {
				dpIn[c] = 0
			}
		}
		if !hasBottom {
			fmt.Fprintln(out, 0)
			continue
		}

		radSame := d
		tmp := int64(d)*int64(d) - 1
		if tmp < 0 {
			tmp = 0
		}
		radCross := int(isqrt(tmp))

		printed := false

		for r := n - 1; r >= 0; r-- {
			row := grid[r]
			nonZero := false
			for c := 0; c < m; c++ {
				if row[c] != 'X' {
					dpIn[c] = 0
				} else if dpIn[c] != 0 {
					nonZero = true
				}
			}
			if !nonZero {
				fmt.Fprintln(out, 0)
				printed = true
				break
			}

			prefix[0] = 0
			for c := 0; c < m; c++ {
				prefix[c+1] = prefix[c] + dpIn[c]
				if prefix[c+1] >= mod {
					prefix[c+1] -= mod
				}
			}
			for c := 0; c < m; c++ {
				if row[c] == 'X' {
					l := c - radSame
					if l < 0 {
						l = 0
					}
					rgt := c + radSame + 1
					if rgt > m {
						rgt = m
					}
					dpLeave[c] = rangeSum(prefix, l, rgt)
				} else {
					dpLeave[c] = 0
				}
			}

			if r == 0 {
				ans := int64(0)
				for c := 0; c < m; c++ {
					ans += dpLeave[c]
				}
				ans %= mod
				fmt.Fprintln(out, ans)
				printed = true
				break
			}

			prefix[0] = 0
			for c := 0; c < m; c++ {
				prefix[c+1] = prefix[c] + dpLeave[c]
				if prefix[c+1] >= mod {
					prefix[c+1] -= mod
				}
			}

			nextRow := grid[r-1]
			nonZero = false
			for c := 0; c < m; c++ {
				if nextRow[c] == 'X' {
					l := c - radCross
					if l < 0 {
						l = 0
					}
					rgt := c + radCross + 1
					if rgt > m {
						rgt = m
					}
					dpIn[c] = rangeSum(prefix, l, rgt)
					if dpIn[c] != 0 {
						nonZero = true
					}
				} else {
					dpIn[c] = 0
				}
			}
			if !nonZero {
				fmt.Fprintln(out, 0)
				printed = true
				break
			}
		}

		if !printed {
			fmt.Fprintln(out, 0)
		}
	}
}
