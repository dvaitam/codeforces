package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(reader, &s)
		n := len(s)
		digits := make([]int, n+1)
		for i := 1; i <= n; i++ {
			digits[i] = int(s[i-1] - '0')
		}
		// prefix sums of digits mod 9
		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = (pref[i-1] + digits[i]) % 9
		}
		var w, m int
		fmt.Fscan(reader, &w, &m)
		// compute value mod9 for each substring of length w
		vals := make([][]int, 9)
		for i := 1; i+w-1 <= n; i++ {
			v := (pref[i+w-1] - pref[i-1]) % 9
			if v < 0 {
				v += 9
			}
			if len(vals[v]) < 2 {
				vals[v] = append(vals[v], i)
			}
		}
		for ; m > 0; m-- {
			var l, r, k int
			fmt.Fscan(reader, &l, &r, &k)
			// value of substring l..r mod 9
			cur := (pref[r] - pref[l-1]) % 9
			if cur < 0 {
				cur += 9
			}
			ans1, ans2 := -1, -1
			// try all possible a,b
			for a := 0; a < 9; a++ {
				if len(vals[a]) == 0 {
					continue
				}
				for b := 0; b < 9; b++ {
					if len(vals[b]) == 0 {
						continue
					}
					if (a*cur+b)%9 != k {
						continue
					}
					if a == b {
						if len(vals[a]) < 2 {
							continue
						}
						p := vals[a][0]
						q := vals[a][1]
						if ans1 == -1 || p < ans1 || (p == ans1 && q < ans2) {
							ans1, ans2 = p, q
						}
					} else {
						p := vals[a][0]
						q := vals[b][0]
						if ans1 == -1 || p < ans1 || (p == ans1 && q < ans2) {
							ans1, ans2 = p, q
						}
					}
				}
			}
			if ans1 == -1 {
				fmt.Fprintln(writer, -1, -1)
			} else {
				fmt.Fprintln(writer, ans1, ans2)
			}
		}
	}
}
