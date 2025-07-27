package main

import (
	"bufio"
	"fmt"
	"os"
)

func buildPrimes(limit int) []int {
	isComp := make([]bool, limit+1)
	primes := []int{}
	for i := 2; i <= limit; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			if i*i <= limit {
				for j := i * i; j <= limit; j += i {
					isComp[j] = true
				}
			}
		}
	}
	return primes
}

func canonical(x int, primes []int) int {
	res := 1
	for _, p := range primes {
		if p*p > x {
			break
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt ^= 1
		}
		if cnt == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	primes := buildPrimes(3200)

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(reader, &a[i])
			a[i] = canonical(a[i], primes)
		}

		// compress values to ids
		idMap := make(map[int]int)
		ids := make([]int, n)
		idCnt := 0
		for i, v := range a {
			if idx, ok := idMap[v]; ok {
				ids[i] = idx
			} else {
				idMap[v] = idCnt
				ids[i] = idCnt
				idCnt++
			}
		}

		m := idCnt
		// precompute next positions
		nxt := make([][]int, n)
		for i := range nxt {
			nxt[i] = make([]int, k+1)
		}
		freq := make([][]int, k+1)
		for i := 0; i <= k; i++ {
			freq[i] = make([]int, m)
		}
		r := make([]int, k+1)
		dup := make([]int, k+1)

		for i := 0; i < n; i++ {
			val := ids[i]
			for j := 0; j <= k; j++ {
				for r[j] < n {
					id := ids[r[j]]
					extra := 0
					if freq[j][id] > 0 {
						extra = 1
					}
					if dup[j]+extra > j {
						break
					}
					dup[j] += extra
					freq[j][id]++
					r[j]++
				}
				nxt[i][j] = r[j]
				// remove current value from window
				if freq[j][val] > 1 {
					dup[j]--
				}
				freq[j][val]--
			}
		}

		const INF = int(1e9)
		dp := make([][]int, n+1)
		for i := 0; i <= n; i++ {
			dp[i] = make([]int, k+1)
			for j := 0; j <= k; j++ {
				dp[i][j] = INF
			}
		}
		dp[0][0] = 0
		for i := 0; i < n; i++ {
			for used := 0; used <= k; used++ {
				if dp[i][used] == INF {
					continue
				}
				for add := 0; add <= k-used; add++ {
					to := nxt[i][add]
					if dp[to][used+add] > dp[i][used]+1 {
						dp[to][used+add] = dp[i][used] + 1
					}
				}
			}
		}
		ans := INF
		for j := 0; j <= k; j++ {
			if dp[n][j] < ans {
				ans = dp[n][j]
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
