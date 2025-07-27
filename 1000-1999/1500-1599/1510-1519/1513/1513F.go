package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	base := int64(0)
	for i := 0; i < n; i++ {
		base += absInt64(a[i] - b[i])
	}

	type pair struct {
		val int64
		idx int
	}
	arrA := make([]pair, n)
	arrB := make([]pair, n)
	for i := 0; i < n; i++ {
		arrA[i] = pair{a[i], i}
		arrB[i] = pair{b[i], i}
	}
	sort.Slice(arrA, func(i, j int) bool { return arrA[i].val < arrA[j].val })
	sort.Slice(arrB, func(i, j int) bool { return arrB[i].val < arrB[j].val })

	aval := make([]int64, n)
	bval := make([]int64, n)
	for i := 0; i < n; i++ {
		aval[i] = arrA[i].val
		bval[i] = arrB[i].val
	}

	best := base
	check := func(i, j int) {
		if i == j {
			return
		}
		cur := base - absInt64(a[i]-b[i]) - absInt64(a[j]-b[j]) +
			absInt64(a[i]-b[j]) + absInt64(a[j]-b[i])
		if cur < best {
			best = cur
		}
	}

	for i := 0; i < n; i++ {
		ai := a[i]
		pos := sort.Search(len(bval), func(k int) bool { return bval[k] >= ai })
		for d := -3; d <= 3; d++ {
			k := pos + d
			if k >= 0 && k < n {
				check(i, arrB[k].idx)
			}
		}
	}

	for j := 0; j < n; j++ {
		bj := b[j]
		pos := sort.Search(len(aval), func(k int) bool { return aval[k] >= bj })
		for d := -3; d <= 3; d++ {
			k := pos + d
			if k >= 0 && k < n {
				check(arrA[k].idx, j)
			}
		}
	}

	fmt.Fprintln(out, best)
}
