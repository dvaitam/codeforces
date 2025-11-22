package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxLog = 40

var pow3 [maxLog]int64

func init() {
	pow3[0] = 1
	for i := 1; i < maxLog; i++ {
		pow3[i] = pow3[i-1] * 3
	}
}

func digitSumBase3(x int64) int64 {
	var s int64
	for x > 0 {
		s += x % 3
		x /= 3
	}
	return s
}

func solveCase(n, k int64) int64 {
	minDeals := digitSumBase3(n)
	if k < minDeals {
		return -1
	}
	if k >= n {
		return 3 * n
	}

	var cnt [maxLog]int64
	cnt[0] = n
	coins := n // current number of deals
	cost := 3 * n

	for i := 0; i+1 < maxLog && coins > k; i++ {
		needMerges := (coins - k + 1) / 2 // merges needed to bring coins to <=k
		avail := cnt[i] / 3               // available merges at this level
		if avail == 0 {
			continue
		}
		t := avail
		if t > needMerges {
			t = needMerges
		}

		cnt[i] -= 3 * t
		cnt[i+1] += t
		coins -= 2 * t
		cost += t * pow3[i] // each merge adds 3^i to total cost
	}

	return cost
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int64
		fmt.Fscan(in, &n, &k)
		fmt.Fprintln(out, solveCase(n, k))
	}
}
