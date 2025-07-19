package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, n1, n2 int
	fmt.Fscan(reader, &n, &n1, &n2)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)
	// First group size is the smaller of n1 and n2
	t := n1
	if n2 < t {
		t = n2
	}
	// Sum largest t values
	sum1 := 0
	for i := n - t; i < n; i++ {
		sum1 += a[i]
	}
	avg1 := float64(sum1) / float64(t)
	// Second group size is the larger of n1 and n2
	t1 := n1
	if n2 > t1 {
		t1 = n2
	}
	// Sum next largest t1 values
	sum2 := 0
	for i := n - t - t1; i < n-t; i++ {
		sum2 += a[i]
	}
	avg2 := float64(sum2) / float64(t1)
	ans := avg1 + avg2
	fmt.Printf("%.6f\n", ans)
}
