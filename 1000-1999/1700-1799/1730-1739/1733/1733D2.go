package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var x, y int64
		fmt.Fscan(in, &n, &x, &y)
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)
		var pos []int
		for i := 0; i < n; i++ {
			if a[i] != b[i] {
				pos = append(pos, i)
			}
		}
		k := len(pos)
		if k%2 == 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		if x >= y {
			if k == 2 && pos[0]+1 == pos[1] {
				if n == 2 || n == 3 {
					fmt.Fprintln(out, x)
				} else if n == 4 && pos[0] == 1 {
					if x < 3*y {
						fmt.Fprintln(out, x)
					} else {
						fmt.Fprintln(out, 3*y)
					}
				} else {
					if x < 2*y {
						fmt.Fprintln(out, x)
					} else {
						fmt.Fprintln(out, 2*y)
					}
				}
			} else {
				fmt.Fprintln(out, int64(k/2)*y)
			}
			continue
		}
		if k == 2 {
			d := pos[1] - pos[0]
			if d == 1 {
				if n == 2 || n == 3 && pos[0] == 1 {
					fmt.Fprintln(out, x)
				} else if n == 4 && pos[0] == 1 {
					if x < 3*y {
						fmt.Fprintln(out, x)
					} else {
						fmt.Fprintln(out, 3*y)
					}
				} else {
					if x < 2*y {
						fmt.Fprintln(out, x)
					} else {
						fmt.Fprintln(out, 2*y)
					}
				}
			} else {
				cost := int64(d) * x
				if y < cost {
					cost = y
				}
				fmt.Fprintln(out, cost)
			}
			continue
		}
		// heuristic for many mismatches
		adj := 0
		for i := 0; i+1 < k; {
			if pos[i+1] == pos[i]+1 {
				adj++
				i += 2
			} else {
				i++
			}
		}
		cost := int64(adj)*x + int64(k-2*adj)/2*y
		fmt.Fprintln(out, cost)
	}
}
