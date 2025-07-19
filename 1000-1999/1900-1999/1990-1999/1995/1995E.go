package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 2000000001

func M(x, n int) int {
	return (x + 2*n) % (2 * n)
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	v := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		fmt.Fscan(in, &v[i])
	}
	if n%2 == 0 {
		maxx := 0
		minn := INF
		for i := 0; i < n/2; i++ {
			s := []int{v[2*i] + v[2*i+1], v[2*i] + v[2*i+n+1], v[2*i+n] + v[2*i+n+1], v[2*i+n] + v[2*i+1]}
			sort.Ints(s)
			if s[2] > maxx {
				maxx = s[2]
			}
			if s[1] < minn {
				minn = s[1]
			}
		}
		fmt.Fprintln(out, maxx-minn)
		return
	}
	if n == 1 {
		fmt.Fprintln(out, 0)
		return
	}
	r := make([]int, 0, 2*n)
	cnt := 0
	for i := 0; i < n; i++ {
		r = append(r, v[cnt])
		cnt ^= 1
		r = append(r, v[cnt])
		cnt = M(cnt+n, n)
	}
	ans := INF
	for id := 0; id < n; id++ {
		for m1 := 0; m1 < 2; m1++ {
			for m2 := 0; m2 < 2; m2++ {
				minn := r[M(2*id-m1, n)] + r[M(2*id+m2+1, n)]
				dp := [2][]int{make([]int, n), make([]int, n)}
				for t := 0; t < 2; t++ {
					for j := 0; j < n; j++ {
						dp[t][j] = INF
					}
				}
				dp[m2][id] = minn
				for j := 1; j < n; j++ {
					d2 := (id + j) % n
					d1 := (id + j - 1) % n
					for c1 := 0; c1 < 2; c1++ {
						for c2 := 0; c2 < 2; c2++ {
							if dp[c1][d1] != INF {
								val := r[M(2*d2-c1, n)] + r[M(2*d2+c2+1, n)]
								if val >= minn {
									cur := dp[c1][d1]
									if val > cur {
										cur = val
									}
									if cur < dp[c2][d2] {
										dp[c2][d2] = cur
									}
								}
							}
						}
					}
				}
				p := (id + n - 1) % n
				if dp[m1][p] != INF {
					diff := dp[m1][p] - minn
					if diff < ans {
						ans = diff
					}
				}
			}
		}
	}
	fmt.Fprintln(out, ans)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solve(in, out)
	}
}
