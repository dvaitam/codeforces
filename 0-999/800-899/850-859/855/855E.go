package main

import (
	"bufio"
	"fmt"
	"os"
)

var dp [11][][]int64
var maxLen [11]int

func digitsLen(b int) int {
	x := int64(1)
	l := 0
	for x <= 1_000_000_000_000_000_000 {
		l++
		x *= int64(b)
	}
	return l
}

func precompute(b int) {
	maxLen[b] = digitsLen(b)
	states := 1 << b
	dp[b] = make([][]int64, maxLen[b]+1)
	dp[b][0] = make([]int64, states)
	dp[b][0][0] = 1
	for l := 1; l <= maxLen[b]; l++ {
		dp[b][l] = make([]int64, states)
		for mask := 0; mask < states; mask++ {
			var sum int64
			for d := 0; d < b; d++ {
				sum += dp[b][l-1][mask^(1<<d)]
			}
			dp[b][l][mask] = sum
		}
	}
}

func digitsBase(x int64, b int) []int {
	if x == 0 {
		return []int{0}
	}
	tmp := make([]int, 0)
	for x > 0 {
		tmp = append(tmp, int(x%int64(b)))
		x /= int64(b)
	}
	for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}
	return tmp
}

func countUpTo(x int64, b int) int64 {
	if x <= 0 {
		return 0
	}
	digits := digitsBase(x, b)
	L := len(digits)
	var res int64
	for l := 1; l < L; l++ {
		for d := 1; d < b; d++ {
			res += dp[b][l-1][1<<d]
		}
	}
	mask := 0
	for i := 0; i < L; i++ {
		cur := digits[i]
		for d := 0; d < cur; d++ {
			if i == 0 && d == 0 {
				continue
			}
			newMask := mask ^ (1 << d)
			remain := L - i - 1
			res += dp[b][remain][newMask]
		}
		mask ^= 1 << cur
	}
	if mask == 0 {
		res++
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var q int
	fmt.Fscan(in, &q)
	for b := 2; b <= 10; b++ {
		precompute(b)
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; q > 0; q-- {
		var b int
		var l, r int64
		fmt.Fscan(in, &b, &l, &r)
		ans := countUpTo(r, b) - countUpTo(l-1, b)
		fmt.Fprintln(out, ans)
	}
}
