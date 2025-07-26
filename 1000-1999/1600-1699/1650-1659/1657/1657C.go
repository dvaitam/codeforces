package main

import (
	"bufio"
	"fmt"
	"os"
)

const base uint64 = 911382323

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		b := []byte(s)

		pow := make([]uint64, n+1)
		pref := make([]uint64, n+1)
		pow[0] = 1
		for i := 0; i < n; i++ {
			pref[i+1] = pref[i]*base + uint64(b[i])
			pow[i+1] = pow[i] * base
		}
		rev := make([]byte, n)
		for i := 0; i < n; i++ {
			rev[i] = b[n-1-i]
		}
		rpref := make([]uint64, n+1)
		for i := 0; i < n; i++ {
			rpref[i+1] = rpref[i]*base + uint64(rev[i])
		}

		isPal := func(l, r int) bool {
			hf := pref[r+1] - pref[l]*pow[r-l+1]
			hr := rpref[n-l] - rpref[n-1-r]*pow[r-l+1]
			return hf == hr
		}

		i := 0
		ops := 0
		for i < n {
			bal := 0
			minBal := 0
			found := false
			for j := i; j < n; j++ {
				if b[j] == '(' {
					bal++
				} else {
					bal--
				}
				if bal < minBal {
					minBal = bal
				}
				if bal == 0 && minBal >= 0 {
					ops++
					i = j + 1
					found = true
					break
				}
				if j-i+1 >= 2 && isPal(i, j) {
					ops++
					i = j + 1
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
		fmt.Fprintln(out, ops, n-i)
	}
}
