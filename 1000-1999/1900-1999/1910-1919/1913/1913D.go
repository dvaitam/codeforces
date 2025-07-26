package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func countReachable(p []int) int64 {
	n := len(p)
	prefix := make([]bool, n)
	mn := int(^uint(0) >> 1) // max int
	for i, x := range p {
		if x < mn {
			mn = x
			prefix[i] = true
		}
	}
	suffix := make([]int, 0)
	mn = int(^uint(0) >> 1)
	for i := n - 1; i >= 0; i-- {
		if p[i] < mn {
			mn = p[i]
			suffix = append(suffix, i)
		}
	}
	stack := make([]int, 0)
	ancSum := make([]int64, 0)
	dp := make([]int64, n)
	sub := make([]int64, n)

	for i, x := range p {
		base := int64(0)
		if prefix[i] {
			base = 1
		}
		poppedTotal := int64(0)
		for len(stack) > 0 && p[stack[len(stack)-1]] > x {
			idx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			ancSum = ancSum[:len(ancSum)-1]
			poppedTotal = (poppedTotal + sub[idx]) % MOD
			sub[i] = (sub[i] + sub[idx]) % MOD
		}
		anc := int64(0)
		if len(ancSum) > 0 {
			anc = ancSum[len(ancSum)-1]
		}
		dp[i] = (base + poppedTotal + anc) % MOD
		sub[i] = (sub[i] + dp[i]) % MOD
		stack = append(stack, i)
		if len(ancSum) > 0 {
			ancSum = append(ancSum, (ancSum[len(ancSum)-1]+dp[i])%MOD)
		} else {
			ancSum = append(ancSum, dp[i])
		}
	}

	res := int64(0)
	for _, idx := range suffix {
		res = (res + dp[idx]) % MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		fmt.Fprintln(out, countReachable(p))
	}
}
