package main

import (
	"bufio"
	"fmt"
	"os"
)

func parity(x int) int {
	x %= 2
	if x < 0 {
		x += 2
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &grid[i])
		}

		minS, maxS := int(1<<60), int(-1<<60)
		minT, maxT := int(1<<60), int(-1<<60)

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 'B' {
					x, y := i+1, j+1
					s := x + y
					tval := x - y
					if s < minS {
						minS = s
					}
					if s > maxS {
						maxS = s
					}
					if tval < minT {
						minT = tval
					}
					if tval > maxT {
						maxT = tval
					}
				}
			}
		}

		check := func(d int) (int, int, bool) {
			L1, R1 := maxS-d, minS+d
			L2, R2 := maxT-d, minT+d
			if L1 > R1 || L2 > R2 {
				return 0, 0, false
			}
			if L1 < 2 {
				L1 = 2
			}
			if R1 > n+m {
				R1 = n + m
			}
			if L2 < 1-m {
				L2 = 1 - m
			}
			if R2 > n-1 {
				R2 = n - 1
			}
			if L1 > R1 || L2 > R2 {
				return 0, 0, false
			}
			for p := L1; p <= L1+1 && p <= R1; p++ {
				for q := L2; q <= L2+1 && q <= R2; q++ {
					if parity(p) != parity(q) {
						continue
					}
					a := (p + q) / 2
					b := (p - q) / 2
					if a >= 1 && a <= n && b >= 1 && b <= m {
						return a, b, true
					}
				}
			}
			return 0, 0, false
		}

		l, r := 0, n+m
		ansA, ansB := 1, 1
		for l <= r {
			mid := (l + r) / 2
			if a, b, ok := check(mid); ok {
				ansA, ansB = a, b
				r = mid - 1
			} else {
				l = mid + 1
			}
		}
		fmt.Fprintf(writer, "%d %d\n", ansA, ansB)
	}
}
