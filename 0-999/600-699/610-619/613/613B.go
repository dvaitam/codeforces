package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type node struct {
	val int64
	idx int
}

func maxBase(arr []node, prefix []int64, length int, budget int64, hi int64) int64 {
	if length == 0 {
		return hi
	}
	lo := arr[0].val
	if lo > hi {
		return lo
	}
	best := lo
	for lo <= hi {
		mid := (lo + hi) / 2
		idx := sort.Search(length, func(i int) bool { return arr[i].val > mid })
		cost := mid*int64(idx) - prefix[idx]
		if cost <= budget {
			best = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var A, cf, cm, m int64
	fmt.Fscan(reader, &n, &A, &cf, &cm, &m)
	arr := make([]node, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i].val)
		arr[i].idx = i
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + arr[i].val
	}
	costA := make([]int64, n+1)
	for k := 1; k <= n; k++ {
		costA[k] = costA[k-1] + (A - arr[n-k].val)
	}
	bestForce := int64(-1)
	bestK := 0
	bestBase := int64(0)
	for k := 0; k <= n; k++ {
		if costA[k] > m {
			continue
		}
		r := m - costA[k]
		base := A
		if k < n {
			base = maxBase(arr, prefix, n-k, r, A-1)
		}
		force := int64(k)*cf + base*cm
		if force > bestForce {
			bestForce = force
			bestK = k
			bestBase = base
		}
	}

	result := make([]int64, n)
	for i := 0; i < n; i++ {
		result[arr[i].idx] = arr[i].val
	}
	for i := 0; i < n-bestK; i++ {
		if result[arr[i].idx] < bestBase {
			result[arr[i].idx] = bestBase
		}
	}
	for i := n - bestK; i < n; i++ {
		result[arr[i].idx] = A
	}
	fmt.Fprintln(writer, bestForce)
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, result[i])
	}
	writer.WriteByte('\n')
}
