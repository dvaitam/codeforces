package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	const maxCapacity = 2000005
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		b := scanner.Bytes()
		x := 0
		for _, c := range b {
			x = x*10 + int(c-'0')
		}
		return x
	}

	if !scanner.Scan() {
		return
	}
	nBytes := scanner.Bytes()
	n := 0
	for _, c := range nBytes {
		n = n*10 + int(c-'0')
	}

	k := scanInt()

	scanner.Scan()
	s := scanner.Bytes()

	if n < 2*k {
		fmt.Println(0)
		return
	}

	const MOD = 1000000007

	wCount := make([]int, n+1)
	bCount := make([]int, n+1)
	xCount := make([]int, n+1)

	for i := 0; i < n; i++ {
		wCount[i+1] = wCount[i]
		bCount[i+1] = bCount[i]
		xCount[i+1] = xCount[i]
		if s[i] == 'W' {
			wCount[i+1]++
		} else if s[i] == 'B' {
			bCount[i+1]++
		} else {
			xCount[i+1]++
		}
	}

	canBeB := func(l, r int) bool {
		return wCount[r+1]-wCount[l] == 0
	}
	canBeW := func(l, r int) bool {
		return bCount[r+1]-bCount[l] == 0
	}

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % MOD
	}

	dpNoB := make([]int64, n+1)
	F := make([]int64, n+1)
	dpNoB[0] = 1

	for i := 1; i <= n; i++ {
		mult := int64(1)
		if s[i-1] == 'X' {
			mult = 2
		}
		val := (dpNoB[i-1] * mult) % MOD

		if i >= k {
			if canBeB(i-k, i-1) {
				if i == k {
					F[i] = 1
				} else {
					if canBeW(i-k-1, i-k-1) {
						F[i] = dpNoB[i-k-1]
					}
				}
			}
		}

		dpNoB[i] = (val - F[i] + MOD) % MOD
	}

	G := make([]int64, n+2)
	G[n] = 1

	for i := n - 1; i >= 0; i-- {
		mult := int64(1)
		if s[i] == 'X' {
			mult = 2
		}
		val := (G[i+1] * mult) % MOD

		sub := int64(0)
		if i+k <= n {
			if canBeW(i, i+k-1) {
				if i+k == n {
					sub = 1
				} else {
					if canBeB(i+k, i+k) {
						sub = G[i+k+1]
					}
				}
			}
		}
		G[i] = (val - sub + MOD) % MOD
	}

	ans := int64(0)
	for l := k; l <= n-k; l++ {
		if F[l] == 0 {
			continue
		}
		cntX := xCount[n] - xCount[l]
		totalWays := pow2[cntX]

		validWays := (totalWays - G[l] + MOD) % MOD
		term := (F[l] * validWays) % MOD
		ans = (ans + term) % MOD
	}

	fmt.Println(ans)
}