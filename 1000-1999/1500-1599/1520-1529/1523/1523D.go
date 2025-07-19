package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"sort"
	"time"
)

const N = 15

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, p int
	if _, err := fmt.Fscan(reader, &n, &m, &p); err != nil {
		return
	}
	a := make([]uint64, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		var x uint64
		for j := 0; j < m; j++ {
			if s[j] == '1' {
				x |= 1 << j
			}
		}
		a[i] = x
	}

	var ans uint64
	dpLen := 1 << N
	dp := make([]int, dpLen)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	trials := 20
	half := (n + 1) / 2
	for t := 0; t < trials; t++ {
		b := a[rnd.Intn(n)]
		// enumerate submasks of b
		submasks := make([]uint64, 0)
		for z := b; ; z = (z - 1) & b {
			submasks = append(submasks, z)
			if z == 0 {
				break
			}
		}
		sort.Slice(submasks, func(i, j int) bool { return submasks[i] < submasks[j] })
		// reset dp
		for i := range dp {
			dp[i] = 0
		}
		// count exact matches
		for _, x := range a {
			xAnd := x & b
			k := sort.Search(len(submasks), func(i int) bool { return submasks[i] >= xAnd })
			dp[k]++
		}
		// SOS DP
		for i := 0; i < N; i++ {
			for j := 0; j < dpLen; j++ {
				if (j>>i)&1 == 1 {
					dp[j^(1<<i)] += dp[j]
				}
			}
		}
		// update answer
		for i, cnt := range dp {
			if cnt >= half {
				if bits.OnesCount64(submasks[i]) > bits.OnesCount64(ans) {
					ans = submasks[i]
				}
			}
		}
	}
	// output result
	for i := 0; i < m; i++ {
		if (ans>>i)&1 == 1 {
			writer.WriteByte('1')
		} else {
			writer.WriteByte('0')
		}
	}
	writer.WriteByte('\n')
}
