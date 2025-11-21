package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func buildFacts(n int) ([]int, []int) {
	fact := make([]int, n+1)
	invFact := make([]int, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * i % mod
	}
	return fact, invFact
}

func comb(n, r int, fact, invFact []int) int {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * invFact[r] % mod * invFact[n-r] % mod
}

func flipBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	for i, ch := range src {
		switch ch {
		case '0':
			dst[i] = '1'
		case '1':
			dst[i] = '0'
		default:
			dst[i] = '?'
		}
	}
	return dst
}

func countCase(s []byte, k int, fact, invFact []int) int {
	n := len(s)
	L := k - 1
	if s[0] == '0' {
		return 0
	}
	// prefixZero[i] — number of forced zeros among first i characters
	prefixZero := make([]int, n+1)
	for i, ch := range s {
		prefixZero[i+1] = prefixZero[i]
		if ch == '0' {
			prefixZero[i+1]++
		}
	}
	// cnt0/cnt1 keep track of forced values for each residue modulo L in the suffix i > l
	cnt0 := make([]int, L)
	cnt1 := make([]int, L)
	state := make([]int, L)
	for i := range state {
		state[i] = -1
	}
	forceZero, forceOne, conflict := 0, 0, 0
	recompute := func(r int) {
		old := state[r]
		switch old {
		case -2:
			conflict--
		case 0:
			forceZero--
		case 1:
			forceOne--
		}
		if cnt0[r] > 0 && cnt1[r] > 0 {
			state[r] = -2
			conflict++
		} else if cnt1[r] > 0 {
			state[r] = 1
			forceOne++
		} else if cnt0[r] > 0 {
			state[r] = 0
			forceZero++
		} else {
			state[r] = -1
		}
	}
	add := func(idx int) {
		ch := s[idx]
		if ch == '?' {
			return
		}
		r := idx % L
		if ch == '1' {
			cnt1[r]++
		} else {
			cnt0[r]++
		}
		recompute(r)
	}
	startL := n - L
	for idx := startL; idx < n; idx++ {
		add(idx)
	}
	target := L / 2
	ans := 0
	// iterate over possible first positions l with d[l,l+L-1]=1
	for l := startL; l >= 1; l-- {
		if prefixZero[l] == 0 && conflict == 0 {
			forcedOnes := forceOne
			forcedZeros := forceZero
			r0 := (l - 1) % L
			addZero := 0
			st := state[r0]
			ok := true
			if l > 1 {
				if st == 1 {
					ok = false
				} else if st == -1 {
					addZero = 1
				}
			}
			if ok {
				forcedZerosAdj := forcedZeros + addZero
				freeAdj := L - forcedOnes - forcedZerosAdj
				needOnes := target - forcedOnes
				if freeAdj >= 0 && needOnes >= 0 && needOnes <= freeAdj && target-forcedZerosAdj == freeAdj-needOnes {
					ans = (ans + comb(freeAdj, needOnes, fact, invFact)) % mod
				}
			}
		}
		if l > 1 {
			add(l - 1)
		}
	}
	// handle the case when every window has difference > 1
	if prefixZero[n-L] == 0 {
		start := n - L
		fixedOnes, fixedZeros, free := 0, 0, 0
		for idx := start; idx < n; idx++ {
			switch s[idx] {
			case '1':
				fixedOnes++
			case '0':
				fixedZeros++
			default:
				free++
			}
		}
		need := target + 1 - fixedOnes
		if need < 0 {
			need = 0
		}
		if need <= free {
			sum := 0
			for choose := need; choose <= free; choose++ {
				sum = (sum + comb(free, choose, fact, invFact)) % mod
			}
			ans = (ans + sum) % mod
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const maxN = 100000 + 5
	fact, invFact := buildFacts(maxN)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		var str string
		fmt.Fscan(in, &n, &k)
		fmt.Fscan(in, &str)
		orig := []byte(str)
		ans := countCase(orig, k, fact, invFact)
		flipped := flipBytes(orig)
		ans = (ans + countCase(flipped, k, fact, invFact)) % mod
		fmt.Fprintln(out, ans)
	}
}
