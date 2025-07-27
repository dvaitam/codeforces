package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	s := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}

	net := make([]int, n)
	minPref := make([]int, n)
	freq := make([]map[int]int, n)
	for i := 0; i < n; i++ {
		str := s[i]
		sum := 0
		minv := 0
		freq[i] = make(map[int]int)
		for _, ch := range str {
			if ch == '(' {
				sum++
			} else {
				sum--
			}
			if sum < minv {
				minv = sum
			}
			if sum == minv {
				freq[i][sum]++
			}
		}
		net[i] = sum
		minPref[i] = minv
	}

	size := 1 << n
	bal := make([]int, size)
	for mask := 1; mask < size; mask++ {
		lb := mask & -mask
		idx := 0
		for (lb>>idx)&1 == 0 {
			idx++
		}
		bal[mask] = bal[mask^lb] + net[idx]
	}

	const negInf = -1 << 60
	dp := make([]int, size)
	for i := range dp {
		dp[i] = negInf
	}
	dp[0] = 0
	ans := 0
	for mask := 0; mask < size; mask++ {
		if dp[mask] < 0 {
			continue
		}
		cur := bal[mask]
		if dp[mask] > ans {
			ans = dp[mask]
		}
		for i := 0; i < n; i++ {
			if mask>>i&1 == 0 {
				add := freq[i][-cur]
				if cur+minPref[i] >= 0 {
					nmask := mask | (1 << i)
					if dp[nmask] < dp[mask]+add {
						dp[nmask] = dp[mask] + add
					}
				} else {
					if dp[mask]+add > ans {
						ans = dp[mask] + add
					}
				}
			}
		}
	}
	fmt.Println(ans)
}
