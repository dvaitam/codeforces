package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		diff := make([]int, n+2)
		for i := 0; i < n; i++ {
			A := a[i] - 1
			B := b[i] - 1
			if A <= B {
				continue
			}
			l := 1
			for l <= A && l <= n {
				q := A / l
				if q == 0 {
					break
				}
				r := A / q
				if r > n {
					r = n
				}
				t0 := B / q
				L := l
				if t0+1 > L {
					L = t0 + 1
				}
				if L <= r {
					diff[L]++
					diff[r+1]--
				}
				l = r + 1
			}
		}

		prefix := 0
		ans := make([]int, 0)
		for k := 1; k <= n; k++ {
			prefix += diff[k]
			if prefix == 0 {
				ans = append(ans, k)
			}
		}

		fmt.Fprintln(out, len(ans))
		if len(ans) > 0 {
			s := make([]string, len(ans))
			for i, v := range ans {
				s[i] = strconv.Itoa(v)
			}
			fmt.Fprintln(out, strings.Join(s, " "))
		} else {
			fmt.Fprintln(out)
		}
	}
}
