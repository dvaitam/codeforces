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
	for T > 0 {
		T--
		var n, k int
		fmt.Fscan(in, &n, &k)
		ch := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			ch[i] = []byte(s)
		}
		// find symmetric pair
		found := false
		n1, n2 := 0, 1
		for i := 0; i < n && !found; i++ {
			for j := i + 1; j < n; j++ {
				if ch[i][j] == ch[j][i] {
					found = true
					n1, n2 = i, j
					break
				}
			}
		}
		// if k is odd or symmetric pair exists
		if k%2 == 1 || found {
			fmt.Fprintln(out, "YES")
			// alternate between n1 and n2
			f := true
			// start at n2
			fmt.Fprintf(out, "%d", n2+1)
			for i := 0; i < k; i++ {
				if f {
					fmt.Fprintf(out, " %d", n1+1)
				} else {
					fmt.Fprintf(out, " %d", n2+1)
				}
				f = !f
			}
			fmt.Fprintln(out)
			continue
		}
		// try to find 3-cycle pattern
		var apc []int
	outer:
		for i := 0; i < n; i++ {
			a, b := -1, -1
			for j := 0; j < n; j++ {
				if j == i {
					continue
				}
				if ch[i][j] == 'a' {
					a = j
				}
				if ch[i][j] == 'b' {
					b = j
				}
			}
			for j := 0; j < n; j++ {
				if j == i {
					continue
				}
				if a >= 0 && ch[j][i] == 'a' {
					apc = []int{j, i, a}
					break outer
				}
				if b >= 0 && ch[j][i] == 'b' {
					apc = []int{j, i, b}
					break outer
				}
			}
		}
		if len(apc) > 0 {
			fmt.Fprintln(out, "YES")
			m := k / 2
			if m%2 == 1 {
				f := true
				// first half
				fmt.Fprintf(out, "%d ", apc[0]+1)
				for t := m; t > 0; t-- {
					if f {
						fmt.Fprintf(out, "%d ", apc[1]+1)
					} else {
						fmt.Fprintf(out, "%d ", apc[0]+1)
					}
					f = !f
				}
				// second half
				f = false
				for t := m; t > 0; t-- {
					if f {
						fmt.Fprintf(out, "%d ", apc[1]+1)
					} else {
						fmt.Fprintf(out, "%d ", apc[2]+1)
					}
					f = !f
				}
			} else {
				f := true
				// first half
				fmt.Fprintf(out, "%d ", apc[1]+1)
				for t := m; t > 0; t-- {
					if f {
						fmt.Fprintf(out, "%d ", apc[0]+1)
					} else {
						fmt.Fprintf(out, "%d ", apc[1]+1)
					}
					f = !f
				}
				// second half
				f = false
				for t := m; t > 0; t-- {
					if f {
						fmt.Fprintf(out, "%d ", apc[1]+1)
					} else {
						fmt.Fprintf(out, "%d ", apc[2]+1)
					}
					f = !f
				}
			}
			fmt.Fprintln(out)
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
