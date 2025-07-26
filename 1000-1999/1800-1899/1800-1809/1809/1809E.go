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

func simulate(start, T, a, b int, ops []int) int {
	c := start
	d := T - c
	for _, v := range ops {
		if v > 0 {
			move := v
			if move > c {
				move = c
			}
			free := b - d
			if move > free {
				move = free
			}
			c -= move
			d += move
		} else {
			move := -v
			if move > d {
				move = d
			}
			space := a - c
			if move > space {
				move = space
			}
			c += move
			d -= move
		}
	}
	return c
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, a, b int
	if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
		return
	}

	ops := make([]int, n)
	adds := make([]int, n)
	prefix := 0
	minPref := 0
	maxPref := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &ops[i])
		adds[i] = -ops[i]
		prefix += adds[i]
		if prefix < minPref {
			minPref = prefix
		}
		if prefix > maxPref {
			maxPref = prefix
		}
	}
	total := prefix

	res := make([][]int, a+1)
	for i := 0; i <= a; i++ {
		res[i] = make([]int, b+1)
	}

	for T := 0; T <= a+b; T++ {
		L := max(0, T-b)
		U := min(a, T)
		if L > U {
			continue
		}
		resL := simulate(L, T, a, b, ops)
		resU := simulate(U, T, a, b, ops)
		thrLow := L - minPref
		thrHigh := U - maxPref
		for c := L; c <= U; c++ {
			var val int
			if c <= thrLow {
				val = resL
			} else if c >= thrHigh {
				val = resU
			} else {
				val = c + total
				if val < L {
					val = L
				} else if val > U {
					val = U
				}
			}
			d := T - c
			if d >= 0 && d <= b {
				res[c][d] = val
			}
		}
	}

	for i := 0; i <= a; i++ {
		for j := 0; j <= b; j++ {
			if j > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, res[i][j])
		}
		writer.WriteByte('\n')
	}
}
