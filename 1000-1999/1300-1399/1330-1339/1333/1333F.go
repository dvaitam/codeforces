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
	// Linear sieve to compute smallest prime factor
	spf := make([]int, n+1)
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > n {
				break
			}
			spf[i*p] = p
		}
	}
	// compute max prime factor for each number
	mpf := make([]int, n+1)
	mpf[1] = 1
	for i := 2; i <= n; i++ {
		p := spf[i]
		mpf[i] = p
		if i/p > 1 {
			if mpf[i/p] > mpf[i] {
				mpf[i] = mpf[i/p]
			}
		}
	}
	freq := make([]int, n+1)
	for i := 1; i <= n; i++ {
		freq[mpf[i]]++
	}
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + freq[i]
	}
	primePre := make([]int, n+1)
	cnt := 0
	isPrime := make([]bool, n+1)
	for _, p := range primes {
		isPrime[p] = true
	}
	for i := 1; i <= n; i++ {
		if isPrime[i] {
			cnt++
		}
		primePre[i] = cnt
	}
	piN := cnt
	M := make([]int, n+1)
	for g := 1; g <= n; g++ {
		M[g] = prefix[g] + (piN - primePre[g])
	}
	res := make([]int, n+1)
	g := 1
	for k := 2; k <= n; k++ {
		for g <= n && M[g] < k {
			g++
		}
		if g > n {
			res[k] = n
		} else {
			res[k] = g
		}
	}
	out := bufio.NewWriter(os.Stdout)
	for k := 2; k <= n; k++ {
		if k > 2 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[k])
	}
	fmt.Fprintln(out)
	out.Flush()
}
