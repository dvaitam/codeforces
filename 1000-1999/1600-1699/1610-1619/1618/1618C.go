package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	br := bufio.NewReader(os.Stdin)
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	var t int
	_, _ = fmt.Fscan(br, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(br, &n)
		v := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(br, &v[i])
		}
		if n < 2 {
			fmt.Fprintln(bw, 0)
			continue
		}
		gc1 := v[0]
		gc2 := v[1]
		for i := 0; i < n; i++ {
			if i&1 == 1 {
				gc2 = gcd(gc2, v[i])
			} else {
				gc1 = gcd(gc1, v[i])
			}
		}
		if gc1 == gc2 {
			fmt.Fprintln(bw, 0)
			continue
		}
		// Try gc1
		if gc1 > 1 {
			ok := true
			for i := 1; i < n; i += 2 {
				if v[i]%gc1 == 0 {
					ok = false
					break
				}
			}
			if ok {
				fmt.Fprintln(bw, gc1)
				continue
			}
		}
		// Try gc2
		ok2 := true
		for i := 0; i < n; i += 2 {
			if v[i]%gc2 == 0 {
				ok2 = false
				break
			}
		}
		if ok2 {
			fmt.Fprintln(bw, gc2)
		} else {
			fmt.Fprintln(bw, 0)
		}
	}
}
