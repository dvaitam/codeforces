package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	if !scanner.Scan() {
		return
	}
	mStr := scanner.Text()
	m, _ := strconv.Atoi(mStr)
	n := scanInt()

	a := make([]int, m)
	var sumA uint64
	for i := 0; i < m; i++ {
		a[i] = scanInt()
		sumA += uint64(a[i])
	}

	b := make([]int, n)
	var sumB uint64
	for i := 0; i < n; i++ {
		b[i] = scanInt()
		sumB += uint64(b[i])
	}

	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	sort.Sort(sort.Reverse(sort.IntSlice(b)))

	minOps := ^uint64(0)

	var currentSumA uint64
	for k := 1; k <= m; k++ {
		currentSumA += uint64(a[k-1])
		ops := (sumA - currentSumA) + uint64(k)*sumB
		if ops < minOps {
			minOps = ops
		}
	}

	var currentSumB uint64
	for l := 1; l <= n; l++ {
		currentSumB += uint64(b[l-1])
		ops := (sumB - currentSumB) + uint64(l)*sumA
		if ops < minOps {
			minOps = ops
		}
	}

	fmt.Println(minOps)
}