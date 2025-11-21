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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)

		ok := true
		consec := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				consec++
				if consec >= k {
					ok = false
					break
				}
			} else {
				consec = 0
			}
		}

		if !ok {
			fmt.Fprintln(out, "NO")
			continue
		}

		perm := make([]int, n)
		nextVal := 1
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				perm[i] = nextVal
				nextVal++
			}
		}
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				perm[i] = nextVal
				nextVal++
			}
		}

		fmt.Fprintln(out, "YES")
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, perm[i])
		}
		fmt.Fprintln(out)
	}
}
