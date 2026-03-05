package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var scanInt = func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if !scanner.Scan() {
		return
	}
	t, _ := strconv.Atoi(scanner.Text())

	for i := 0; i < t; i++ {
		n := scanInt()
		m := scanInt()

		b := make([]int64, n+1)
		for j := 1; j <= n; j++ {
			b[j] = int64(scanInt())
		}

		type Edge struct {
			to int
			w  int64
		}
		adj := make([][]Edge, n+1)

		for j := 0; j < m; j++ {
			u := scanInt()
			v := scanInt()
			w := int64(scanInt())
			adj[u] = append(adj[u], Edge{v, w})
		}

		dp := make([]int64, n+1)

		check := func(x int64) bool {
			for k := 1; k <= n; k++ {
				dp[k] = -1
			}

			if b[1] >= x {
				dp[1] = x
			} else {
				dp[1] = b[1]
			}

			for u := 1; u < n; u++ {
				if dp[u] == -1 {
					continue
				}
				cur := dp[u]
				for _, e := range adj[u] {
					if e.w <= cur {
						nextBat := cur + b[e.to]
						if nextBat > x {
							nextBat = x
						}
						if nextBat > dp[e.to] {
							dp[e.to] = nextBat
						}
					}
				}
			}
			return dp[n] != -1
		}

		const maxVal int64 = 200000000000000 // 2 * 10^14 max possible answer
		if !check(maxVal) {
			fmt.Fprintln(writer, -1)
		} else {
			low, high := int64(0), maxVal
			ans := int64(-1)
			for low <= high {
				mid := low + (high-low)/2
				if check(mid) {
					ans = mid
					high = mid - 1
				} else {
					low = mid + 1
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
