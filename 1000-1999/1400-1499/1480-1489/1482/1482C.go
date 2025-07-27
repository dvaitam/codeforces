package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for Codeforces problem 1482C - Basic Diplomacy.
// We first pick the first available friend each day. If some
// friend is chosen more than ceil(m/2) times, we reassign some of
// their days to other available friends (when possible) until the
// limit is satisfied or we determine it is impossible.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		days := make([][]int, m)
		for i := 0; i < m; i++ {
			var k int
			fmt.Fscan(in, &k)
			list := make([]int, k)
			for j := 0; j < k; j++ {
				fmt.Fscan(in, &list[j])
			}
			days[i] = list
		}

		limit := (m + 1) / 2
		choose := make([]int, m)
		cnt := make([]int, n+1)
		for i := 0; i < m; i++ {
			choose[i] = days[i][0]
			cnt[choose[i]]++
		}

		bad := 0
		for i := 1; i <= n; i++ {
			if cnt[i] > cnt[bad] {
				bad = i
			}
		}

		if cnt[bad] > limit {
			for i := 0; i < m && cnt[bad] > limit; i++ {
				if choose[i] == bad && len(days[i]) > 1 {
					for _, f := range days[i] {
						if f != bad {
							choose[i] = f
							cnt[bad]--
							break
						}
					}
				}
			}
			if cnt[bad] > limit {
				fmt.Fprintln(out, "NO")
				continue
			}
		}

		fmt.Fprintln(out, "YES")
		for i := 0; i < m; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, choose[i])
		}
		fmt.Fprintln(out)
	}
}
