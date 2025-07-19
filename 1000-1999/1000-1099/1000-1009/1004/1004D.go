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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	cnt := make([]int, t+1)
	Mx := 0
	for i := 0; i < t; i++ {
		var d int
		fmt.Fscan(reader, &d)
		if d <= t {
			cnt[d]++
		}
		if d > Mx {
			Mx = d
		}
	}
	Bl := 0
	for i := 1; i <= t; i++ {
		if cnt[i] != i*4 {
			Bl = i - 1
			break
		}
	}
	// Try possible dimensions
	for n := 1; n*n <= t; n++ {
		if t%n != 0 {
			continue
		}
		m := t / n
		if n+m-2 < Mx {
			continue
		}
		for X := 1; X <= n; X++ {
			Y := Mx - X + 2
			if Y <= 0 || Y > m {
				continue
			}
			// global minimal distance from borders
			g := min(min(X-1, Y-1), min(n-X, m-Y))
			if g != Bl {
				continue
			}
			// verify counts
			c2 := make([]int, t+1)
			for i := 1; i <= n; i++ {
				for j := 1; j <= m; j++ {
					d := abs(i-X) + abs(j-Y)
					if d <= t {
						c2[d]++
					}
				}
			}
			ok := true
			for i := 1; i <= t; i++ {
				if c2[i] != cnt[i] {
					ok = false
					break
				}
			}
			if ok {
				fmt.Fprintf(writer, "%d %d\n%d %d\n", n, m, X, Y)
				return
			}
		}
	}
	fmt.Fprintln(writer, -1)
}
