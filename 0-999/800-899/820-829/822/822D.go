package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAXN int = 5000000

var spf []int
var f []int64

func sieve(n int) {
	spf = make([]int, n+1)
	primes := []int{}
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
}

func precompute(n int) {
	f = make([]int64, n+1)
	f[1] = 0
	for i := 2; i <= n; i++ {
		p := spf[i]
		f[i] = f[i/p] + int64(i)*(int64(p)-1)/2
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int64
	var l, r int
	if _, err := fmt.Fscan(reader, &t, &l, &r); err != nil {
		return
	}

	sieve(r)
	precompute(r)

	res := int64(0)
	pow := int64(1)
	for i := l; i <= r; i++ {
		res = (res + pow*(f[i]%MOD)) % MOD
		pow = pow * t % MOD
	}
	fmt.Fprintln(writer, res%MOD)
}
