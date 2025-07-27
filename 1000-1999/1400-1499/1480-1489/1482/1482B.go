package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(a []int64) (int64, int64, bool) {
	n := len(a)
	if n == 1 {
		return 0, 0, true
	}
	cSet := false
	var c int64
	for i := 1; i < n; i++ {
		diff := a[i] - a[i-1]
		if diff >= 0 {
			if !cSet {
				cSet = true
				c = diff
			} else if c != diff {
				return 0, 0, false
			}
		}
	}
	if !cSet {
		diff0 := a[1] - a[0]
		for i := 2; i < n; i++ {
			if a[i]-a[i-1] != diff0 {
				return 0, 0, false
			}
		}
		return 0, 0, true
	}
	m := int64(-1)
	for i := 1; i < n; i++ {
		diff := a[i] - a[i-1]
		if diff < 0 {
			cand := c - diff
			if m == -1 {
				m = cand
			} else if m != cand {
				return 0, 0, false
			}
		}
	}
	if m == -1 {
		return 0, 0, true
	}
	if c >= m {
		return 0, 0, false
	}
	maxA := a[0]
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
	}
	if maxA >= m {
		return 0, 0, false
	}
	cur := a[0]
	for i := 1; i < n; i++ {
		cur = (cur + c) % m
		if cur != a[i] {
			return 0, 0, false
		}
	}
	return m, c, true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		m, c, ok := solve(arr)
		if !ok {
			fmt.Fprintln(writer, -1)
		} else if m == 0 {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintf(writer, "%d %d\n", m, c)
		}
	}
}
