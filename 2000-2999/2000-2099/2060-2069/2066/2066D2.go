package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

// We process floors in increasing order (1..n). For each k we consider the stage
// during which floor k is still allowed. At the start of stage k exactly "x"
// planes coming from floors >= k have already happened. The stage length is
// len = c - x (may be zero). During this segment only floors >= k are allowed
// and floor k cannot appear afterwards. We count how many of the len events are
// from floors >= k+1 (denote h); then the start value for the next stage becomes
// x+h. Because c <= 100, len never exceeds 100, so we can run a small DP per
// segment to enumerate the number of ways for every possible h that matches the
// template.

type key struct {
	k   int
	x   int
	pos int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; T > 0; T-- {
		var n, c, m int
		fmt.Fscan(in, &n, &c, &m)
		a := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &a[i])
		}

		memo := make(map[key]int)

		var solve func(k, x, pos int) int
		solve = func(k, x, pos int) int {
			if x > c {
				return 0
			}
			if k == n+1 {
				if pos == m {
					return 1
				}
				return 0
			}
			st := key{k, x, pos}
			if val, ok := memo[st]; ok {
				return val
			}
			lenStage := c - x
			if pos+lenStage > m {
				return 0
			}

			// segment DP to compute number of ways to fill the next lenStage
			// positions, grouped by how many are from floors >= k+1 ("high").
			ways := make([]int, lenStage+1)
			tmp := make([]int, lenStage+1)
			ways[0] = 1
			for i := 0; i < lenStage; i++ {
				for j := 0; j <= lenStage; j++ {
					tmp[j] = 0
				}
				val := a[pos+i]
				if val != 0 && val < k {
					// incompatible value
					memo[st] = 0
					return 0
				}
				for h := 0; h <= i; h++ {
					cur := ways[h]
					if cur == 0 {
						continue
					}
					if val == 0 {
						// choose floor k
						tmp[h] = (tmp[h] + cur) % mod
						// choose any floor in [k+1, n]
						highChoices := n - k
						if highChoices > 0 {
							tmp[h+1] = (tmp[h+1] + cur*highChoices) % mod
						}
					} else if val == k {
						tmp[h] = (tmp[h] + cur) % mod
					} else { // val > k
						if val <= n {
							tmp[h+1] = (tmp[h+1] + cur) % mod
						} else {
							memo[st] = 0
							return 0
						}
					}
				}
				ways, tmp = tmp, ways
			}

			res := 0
			for h := 0; h <= lenStage; h++ {
				if ways[h] == 0 {
					continue
				}
				res = (res + ways[h]*solve(k+1, x+h, pos+lenStage)) % mod
			}
			memo[st] = res
			return res
		}

		ans := solve(1, 0, 0)
		fmt.Fprintln(out, ans)
	}
}
