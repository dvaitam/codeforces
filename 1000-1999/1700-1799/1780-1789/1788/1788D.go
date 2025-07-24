package main

// Attempted solution for problem 1788D based on interpreting the statement in
// problemD.txt. The approach counts subsets for which a pair of adjacent dots
// moves apart, which we consider as a boundary resulting in an extra final
// coordinate.  For each pair (j,k) we use binary search over the positions to
// count valid left and right neighbours that satisfy distance constraints.
// The overall complexity is roughly O(n^2 log n).

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1_000_000_007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	x := make([]int, n)
	for i := range x {
		fmt.Fscan(reader, &x[i])
	}

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
	}

	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = (prefix[i-1] + pow2[i-1]) % mod
	}
	suffix := make([]int64, n+2)
	for i := n; i >= 1; i-- {
		suffix[i] = (suffix[i+1] + pow2[n-i]) % mod
	}

	res := int64(0)
	for j := 1; j <= n-3; j++ {
		for k := j + 1; k <= n-2; k++ {
			mid := x[k] - x[j]
			// left neighbours
			leftIdx := sort.Search(j, func(t int) bool {
				return x[t] >= x[j]-mid
			})
			leftSum := (prefix[j] - prefix[leftIdx]) % mod
			if leftSum < 0 {
				leftSum += mod
			}
			// right neighbours
			idx := sort.Search(n-k-1, func(t int) bool {
				return x[k+1+t] >= x[k]+mid
			})
			rightIdx := k + idx - 1
			var rightSum int64
			if rightIdx >= k+1 {
				rightSum = (suffix[k+1] - suffix[rightIdx+1]) % mod
				if rightSum < 0 {
					rightSum += mod
				}
			}
			res = (res + leftSum*rightSum) % mod
		}
	}

	total := (pow2[n] - int64(n) - 1) % mod
	if total < 0 {
		total += mod
	}
	ans := (total + res) % mod
	fmt.Fprintln(writer, ans)
}
