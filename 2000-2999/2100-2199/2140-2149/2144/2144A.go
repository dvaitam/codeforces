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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + a[i-1]
		}
		found := false
		lans, rans := 0, 0
		for l := 1; l <= n-2 && !found; l++ {
			s1 := ((pref[l] % 3) + 3) % 3
			for r := l + 1; r <= n-1; r++ {
				s2 := ((pref[r]-pref[l])%3 + 3) % 3
				s3 := ((pref[n]-pref[r])%3 + 3) % 3
				if (s1 == s2 && s2 == s3) || (s1 != s2 && s1 != s3 && s2 != s3) {
					found = true
					lans = l
					rans = r
					break
				}
			}
		}
		if !found {
			fmt.Fprintln(out, "0 0")
		} else {
			fmt.Fprintln(out, lans, rans)
		}
	}
}
