package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	limit := 1000000
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}

	freq := make(map[int]int)
	for _, x := range arr {
		freq[x]++
	}

	fact := make([]int64, 2*n+1)
	fact[0] = 1
	for i := 1; i <= 2*n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}

	var primeCounts []int
	for v, c := range freq {
		if v <= limit && isPrime[v] {
			primeCounts = append(primeCounts, c)
		}
	}

	dp := make([]int64, n+1)
	dp[0] = 1
	for _, c := range primeCounts {
		for j := n; j >= 1; j-- {
			dp[j] = (dp[j] + dp[j-1]*int64(c)) % mod
		}
	}

	constProd := int64(1)
	Aconst := constProd
	for _, c := range freq {
		Aconst = Aconst * fact[c] % mod
	}

	ans := fact[n] * dp[n] % mod
	ans = ans * modPow(Aconst, mod-2) % mod

	fmt.Fprintln(out, ans)
}
