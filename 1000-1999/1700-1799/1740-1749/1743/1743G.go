package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const BASE uint64 = 911382323

func computeHash(b []byte) uint64 {
	var h uint64
	for _, c := range b {
		h = h*BASE + uint64(c)
	}
	return h
}

func buildFibHashes(limit int) ([]int, []uint64) {
	f0 := []byte{'0'}
	f1 := []byte{'1'}
	lengths := []int{1, 1}
	hashes := []uint64{computeHash(f0), computeHash(f1)}
	for {
		nextLen := len(f1) + len(f0)
		if nextLen > limit {
			break
		}
		next := make([]byte, nextLen)
		copy(next, f1)
		copy(next[len(f1):], f0)
		lengths = append(lengths, nextLen)
		hashes = append(hashes, computeHash(next))
		f0, f1 = f1, next
	}
	return lengths, hashes
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	strs := make([]string, n)
	total := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &strs[i])
		total += len(strs[i])
	}

	lengths, hashes := buildFibHashes(total)

	pow := make([]uint64, total+1)
	prefHash := make([]uint64, total+1)
	dp := make([]int64, total+1)
	prefDP := make([]int64, total+1)

	pow[0] = 1
	dp[0] = 1
	prefDP[0] = 1

	pos := 0
	res := make([]int64, n)
	fibCount := len(lengths)

	for idx, s := range strs {
		for j := 0; j < len(s); j++ {
			c := s[j]
			pos++
			pow[pos] = pow[pos-1] * BASE
			prefHash[pos] = prefHash[pos-1]*BASE + uint64(c)
			val := prefDP[pos-1]
			for k := 0; k < fibCount; k++ {
				L := lengths[k]
				if L > pos {
					break
				}
				start := pos - L
				sub := prefHash[pos] - prefHash[start]*pow[L]
				if sub == hashes[k] {
					val -= dp[start]
				}
			}
			val %= MOD
			if val < 0 {
				val += MOD
			}
			dp[pos] = val
			prefDP[pos] = (prefDP[pos-1] + val) % MOD
		}
		res[idx] = dp[pos]
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, res[i])
	}
	out.WriteByte('\n')
}
