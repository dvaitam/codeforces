package main

import (
	"bufio"
	"fmt"
	"os"
)

func lis(arr []int) int {
	b := make([]int, 0)
	for _, x := range arr {
		i := lowerBound(b, x)
		if i == len(b) {
			b = append(b, x)
		} else {
			b[i] = x
		}
	}
	return len(b)
}

func lowerBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) / 2
		if a[m] < x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}
		best := lis(arr)
		rev := make([]int, n)
		for i := 0; i < n; i++ {
			rev[i] = arr[n-1-i]
		}
		if v := lis(rev); v > best {
			best = v
		}
		fmt.Fprintln(out, best)
	}
}
