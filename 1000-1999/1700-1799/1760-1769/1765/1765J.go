package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}
	combined := make([]int, 0, 2*n)
	combined = append(combined, a...)
	combined = append(combined, b...)
	sort.Ints(combined)
	m1 := combined[n-1]
	m2 := combined[n]
	var cost1, cost2 int64
	for i := 0; i < n; i++ {
		if a[i] >= m1 {
			cost1 += int64(a[i] - m1)
		} else {
			cost1 += int64(m1 - a[i])
		}
		if a[i] >= m2 {
			cost2 += int64(a[i] - m2)
		} else {
			cost2 += int64(m2 - a[i])
		}
	}
	for i := 0; i < n; i++ {
		if b[i] >= m1 {
			cost1 += int64(b[i] - m1)
		} else {
			cost1 += int64(m1 - b[i])
		}
		if b[i] >= m2 {
			cost2 += int64(b[i] - m2)
		} else {
			cost2 += int64(m2 - b[i])
		}
	}
	if cost2 < cost1 {
		cost1 = cost2
	}
	fmt.Println(cost1)
}
