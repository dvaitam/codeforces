package main

import (
	"bufio"
	"fmt"
	"os"
)

func buildSparse(arr []int64) ([][]int64, []int) {
	n := len(arr)
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}
	k := log[n]
	st := make([][]int64, k+1)
	st[0] = make([]int64, n)
	copy(st[0], arr)
	for j := 1; j <= k; j++ {
		length := 1 << j
		st[j] = make([]int64, n-length+1)
		for i := 0; i+length <= n; i++ {
			v1 := st[j-1][i]
			v2 := st[j-1][i+(length>>1)]
			if v1 < v2 {
				st[j][i] = v1
			} else {
				st[j][i] = v2
			}
		}
	}
	return st, log
}

func queryMin(st [][]int64, log []int, l, r int) int64 {
	length := r - l + 1
	k := log[length]
	v1 := st[k][l]
	v2 := st[k][r-(1<<k)+1]
	if v1 < v2 {
		return v1
	}
	return v2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var s int64
		fmt.Fscan(reader, &n, &s)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		prefix := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			prefix[i] = prefix[i-1] + a[i-1]
		}
		st, log := buildSparse(prefix)
		bestLen, bestL, bestR := 0, -1, -1
		for l := 1; l <= n; l++ {
			threshold := prefix[l-1] - s
			lo, hi := l, n
			ans := l - 1
			for lo <= hi {
				mid := (lo + hi) / 2
				if queryMin(st, log, l, mid) >= threshold {
					ans = mid
					lo = mid + 1
				} else {
					hi = mid - 1
				}
			}
			if ans >= l {
				if ans-l+1 > bestLen {
					bestLen = ans - l + 1
					bestL = l
					bestR = ans
				}
			}
		}
		if bestLen == 0 {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, bestL, bestR)
		}
	}
}
