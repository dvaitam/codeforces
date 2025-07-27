package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007
const maxVal int = 200000

var spf [maxVal + 1]int

func initSieve() {
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxVal; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

func factorize(x int) map[int]int {
	res := make(map[int]int)
	for x > 1 {
		p := spf[x]
		if p == 0 {
			p = x
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		res[p] += cnt
	}
	return res
}

type primeData struct {
	indexExp map[int]int
	cnt      map[int]int
	minExp   int
	missing  int
}

func modPow(a, e int) int {
	res := 1
	b := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * b % MOD
		}
		b = b * b % MOD
		e >>= 1
	}
	return res
}

func minKey(m map[int]int) int {
	min := -1
	for k := range m {
		if min == -1 || k < min {
			min = k
		}
	}
	if min == -1 {
		return 0
	}
	return min
}

func main() {
	initSieve()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	primes := make(map[int]*primeData)

	for i, val := range a {
		factors := factorize(val)
		for p, e := range factors {
			pd := primes[p]
			if pd == nil {
				pd = &primeData{indexExp: make(map[int]int), cnt: make(map[int]int), missing: n}
				primes[p] = pd
			}
			pd.indexExp[i] = e
			pd.cnt[e]++
			pd.missing--
		}
	}

	curGCD := 1
	for p, pd := range primes {
		if pd.missing == 0 {
			pd.minExp = minKey(pd.cnt)
			curGCD = curGCD * modPow(p, pd.minExp) % MOD
		}
	}

	for ; q > 0; q-- {
		var idx, x int
		fmt.Fscan(reader, &idx, &x)
		idx--
		factors := factorize(x)
		for p, add := range factors {
			pd := primes[p]
			if pd == nil {
				pd = &primeData{indexExp: make(map[int]int), cnt: make(map[int]int), missing: n}
				primes[p] = pd
			}
			oldExp := pd.indexExp[idx]
			if oldExp == 0 {
				pd.missing--
			} else {
				pd.cnt[oldExp]--
				if pd.cnt[oldExp] == 0 {
					delete(pd.cnt, oldExp)
				}
			}
			newExp := oldExp + add
			pd.indexExp[idx] = newExp
			pd.cnt[newExp]++

			oldMin := pd.minExp
			var newMin int
			if pd.missing > 0 {
				newMin = 0
			} else {
				if (oldExp == oldMin && pd.cnt[oldExp] == 0) || oldMin == 0 {
					newMin = minKey(pd.cnt)
				} else {
					newMin = oldMin
				}
			}
			if newMin > oldMin {
				curGCD = curGCD * modPow(p, newMin-oldMin) % MOD
				pd.minExp = newMin
			} else {
				pd.minExp = newMin
			}
		}
		fmt.Fprintln(writer, curGCD)
	}
}
