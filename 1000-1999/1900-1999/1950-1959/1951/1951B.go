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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		r := a[k-1]

		prefix := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if a[i-1] > prefix[i-1] {
				prefix[i] = a[i-1]
			} else {
				prefix[i] = prefix[i-1]
			}
		}

		prefixWithout := make([]int, n+1)
		mx := 0
		for i := 1; i <= n; i++ {
			if i == k {
				prefixWithout[i] = mx
			} else {
				if a[i-1] > mx {
					mx = a[i-1]
				}
				prefixWithout[i] = mx
			}
		}

		ng := make([]int, n+2)
		nextIdx := n + 1
		for i := n; i >= 1; i-- {
			if a[i-1] > r {
				nextIdx = i
			}
			ng[i] = nextIdx
		}

		ans := 0
		for p := 1; p <= n; p++ {
			b := a[p-1]
			var pref int
			if p <= k {
				pref = prefix[p-1]
			} else {
				pref = prefixWithout[p-1]
				if k <= p-1 && b > pref {
					pref = b
				}
			}

			q := ng[p+1]
			if p < k && b > r && k < q {
				q = k
			}

			var wins int
			if p == 1 {
				if q == n+1 {
					wins = n - 1
				} else {
					wins = q - 2
				}
			} else {
				if r <= pref {
					wins = 0
				} else {
					if q == n+1 {
						wins = n - p + 1
					} else {
						wins = q - p
					}
				}
			}
			if wins > ans {
				ans = wins
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
