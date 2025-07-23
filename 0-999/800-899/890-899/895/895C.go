package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	freq := make([]int, 71)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x >= 1 && x <= 70 {
			freq[x]++
		}
	}

	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67}
	m := len(primes)
	mask := make([]int, 71)
	for v := 1; v <= 70; v++ {
		cur := v
		mm := 0
		for i, p := range primes {
			cnt := 0
			for cur%p == 0 {
				cur /= p
				cnt ^= 1
			}
			if cnt == 1 {
				mm |= 1 << i
			}
		}
		mask[v] = mm
	}

	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % MOD
	}

	size := 1 << m
	dp := make([]int, size)
	dp[0] = 1
	for v := 1; v <= 70; v++ {
		c := freq[v]
		if c == 0 {
			continue
		}
		even := pow2[c-1]
		odd := pow2[c-1]
		msk := mask[v]
		newDP := make([]int, size)
		for j := 0; j < size; j++ {
			if dp[j] == 0 {
				continue
			}
			val := dp[j]
			ne := (val * even) % MOD
			newDP[j] = (newDP[j] + ne) % MOD
			no := (val * odd) % MOD
			j2 := j ^ msk
			newDP[j2] = (newDP[j2] + no) % MOD
		}
		dp = newDP
	}

	ans := dp[0] - 1
	if ans < 0 {
		ans += MOD
	}
	fmt.Fprintln(writer, ans)
}
