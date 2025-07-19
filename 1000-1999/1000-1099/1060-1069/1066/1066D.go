package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	rdr := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(rdr, &n, &m, &k); err != nil {
		return
	}
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &A[i])
	}

	check := func(s int) bool {
		sum := 0
		cnt := 0
		for i := s; i < n; i++ {
			sum += A[i]
			if sum > k {
				cnt++
				sum = A[i]
			}
		}
		cnt++
		return cnt <= m
	}

	l, r := 0, n-1
	for l < r {
		mid := (l + r) / 2
		if check(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	fmt.Fprintln(w, n-l)
}
