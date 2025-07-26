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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			arr[i]--
		}
		// prefix frequency for each value
		pref := make([][26]int, n+1)
		for i := 0; i < n; i++ {
			copy(pref[i+1][:], pref[i][:])
			pref[i+1][arr[i]]++
		}
		// positions for each value
		pos := make([][]int, 26)
		for i, v := range arr {
			pos[v] = append(pos[v], i)
		}
		ans := 0
		for c := 0; c < 26; c++ {
			if len(pos[c]) > ans {
				ans = len(pos[c])
			}
		}
		for c := 0; c < 26; c++ {
			p := pos[c]
			m := len(p)
			for x := 1; x*2 <= m; x++ {
				L := p[x-1] + 1
				R := p[m-x] - 1
				if L > R {
					break
				}
				best := 0
				for b := 0; b < 26; b++ {
					count := pref[R+1][b] - pref[L][b]
					if count > best {
						best = count
					}
				}
				if best+2*x > ans {
					ans = best + 2*x
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
