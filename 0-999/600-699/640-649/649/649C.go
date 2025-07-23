package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var x, y int64
	if _, err := fmt.Fscan(reader, &n, &x, &y); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)

	prefixDouble := make([]int64, n+1)
	prefixOdd := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefixDouble[i+1] = prefixDouble[i] + int64((a[i]+1)/2)
		if a[i]%2 == 1 {
			prefixOdd[i+1] = prefixOdd[i] + 1
		} else {
			prefixOdd[i+1] = prefixOdd[i]
		}
	}

	feasible := func(k int) bool {
		totalDouble := prefixDouble[k]
		need := totalDouble - x
		if need <= 0 {
			// have enough double-sided sheets
			return true
		}
		odd := prefixOdd[k]
		use1 := odd
		if use1 > y {
			use1 = y
		}
		if use1 > need {
			use1 = need
		}
		need -= use1
		yLeft := y - use1
		if yLeft >= need*2 {
			return true
		}
		return false
	}

	l, r := 0, n
	for l < r {
		mid := (l + r + 1) / 2
		if feasible(mid) {
			l = mid
		} else {
			r = mid - 1
		}
	}
	fmt.Fprintln(writer, l)
}
