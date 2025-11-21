package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func prefixDiff(arr []int) []int64 {
	sort.Ints(arr)
	pairs := len(arr) / 2
	res := make([]int64, pairs+1)
	for i := 1; i <= pairs; i++ {
		diff := int64(arr[len(arr)-i] - arr[i-1])
		res[i] = res[i-1] + diff
	}
	return res
}

func max3(a, b, c int) int {
	if a < b {
		a = b
	}
	if a < c {
		a = c
	}
	return a
}

func min3(a, b, c int) int {
	if a > b {
		a = b
	}
	if a > c {
		a = c
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
		}

		preA := prefixDiff(a)
		preB := prefixDiff(b)
		pairsA := len(preA) - 1
		pairsB := len(preB) - 1

		kmax := n + m
		if temp := (n + m) / 3; temp < kmax {
			kmax = temp
		}
		if n < kmax {
			kmax = n
		}
		if m < kmax {
			kmax = m
		}

		fmt.Fprintln(writer, kmax)
		if kmax == 0 {
			continue
		}

		ans := make([]int64, kmax+1)
		for k := 1; k <= kmax; k++ {
			L := max3(0, k-pairsB, 2*k-m)
			R := min3(pairsA, n-k, k)
			if L > R {
				ans[k] = 0
				continue
			}
			l := L
			r := R
			for r-l > 6 {
				m1 := l + (r-l)/3
				m2 := r - (r-l)/3
				val1 := preA[m1] + preB[k-m1]
				val2 := preA[m2] + preB[k-m2]
				if val1 < val2 {
					l = m1 + 1
				} else {
					r = m2 - 1
				}
			}
			best := int64(0)
			for x := l; x <= r; x++ {
				val := preA[x] + preB[k-x]
				if val > best {
					best = val
				}
			}
			ans[k] = best
		}
		for k := 1; k <= kmax; k++ {
			if k > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans[k])
		}
		fmt.Fprintln(writer)
	}
}
