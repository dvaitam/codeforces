package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return
	}
	n, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	k, _ := strconv.Atoi(scanner.Text())

	a := make([]int, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		a[i], _ = strconv.Atoi(scanner.Text())
	}

	sort.Ints(a)

	MOD := 998244353
	MaxV := (a[n-1] - a[0]) / (k - 1)

	ans := 0
	dp := make([]int, n)
	new_dp := make([]int, n)

	for v := 1; v <= MaxV; v++ {
		for i := 0; i < n; i++ {
			dp[i] = 1
		}
		for j := 2; j <= k; j++ {
			sum := 0
			p := 0
			allZero := true
			for i := 0; i < n; i++ {
				for p < i && a[i]-a[p] >= v {
					sum += dp[p]
					if sum >= MOD {
						sum -= MOD
					}
					p++
				}
				new_dp[i] = sum
				if sum > 0 {
					allZero = false
				}
			}
			for i := 0; i < n; i++ {
				dp[i] = new_dp[i]
			}
			if allZero {
				break
			}
		}
		for i := 0; i < n; i++ {
			ans += dp[i]
			if ans >= MOD {
				ans -= MOD
			}
		}
	}
	fmt.Println(ans)
}