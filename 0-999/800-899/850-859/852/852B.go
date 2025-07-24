package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func conv(a, b []int64, modVal int) []int64 {
	res := make([]int64, modVal)
	for i, av := range a {
		if av == 0 {
			continue
		}
		for j, bv := range b {
			if bv == 0 {
				continue
			}
			res[(i+j)%modVal] = (res[(i+j)%modVal] + av*bv) % MOD
		}
	}
	return res
}

func powVec(base []int64, exp int, modVal int) []int64 {
	res := make([]int64, modVal)
	res[0] = 1
	for exp > 0 {
		if exp&1 == 1 {
			res = conv(res, base, modVal)
		}
		exp >>= 1
		if exp > 0 {
			base = conv(base, base, modVal)
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, L, M int
	if _, err := fmt.Fscan(in, &N, &L, &M); err != nil {
		return
	}

	cntA := make([]int64, M)
	cntB := make([]int64, M)
	bResid := make([]int, N)
	for i := 0; i < N; i++ {
		var x int
		fmt.Fscan(in, &x)
		cntA[x%M]++
	}
	for i := 0; i < N; i++ {
		var x int
		fmt.Fscan(in, &x)
		r := x % M
		cntB[r]++
		bResid[i] = r
	}
	lastFreq := make([][]int64, M)
	for i := 0; i < M; i++ {
		lastFreq[i] = make([]int64, M)
	}
	for i := 0; i < N; i++ {
		var x int
		fmt.Fscan(in, &x)
		c := x % M
		b := bResid[i]
		lastFreq[b][c]++
	}

	power := L - 2
	midDist := make([]int64, M)
	if power >= 0 {
		midDist = powVec(cntB, power, M)
	}

	var ans int64
	for a := 0; a < M; a++ {
		if cntA[a] == 0 {
			continue
		}
		for b := 0; b < M; b++ {
			row := lastFreq[b]
			if row == nil {
				continue
			}
			for c := 0; c < M; c++ {
				freq := row[c]
				if freq == 0 {
					continue
				}
				sum := (a + b + c) % M
				need := (M - sum) % M
				ways := midDist[need]
				if ways == 0 {
					continue
				}
				val := cntA[a] % MOD
				val = (val * freq) % MOD
				val = (val * ways) % MOD
				ans = (ans + val) % MOD
			}
		}
	}

	fmt.Fprintln(out, ans%MOD)
}
