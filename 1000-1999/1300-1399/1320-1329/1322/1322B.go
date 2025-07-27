package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	ans := 0
	// maximum bit needed is 24 since ai <= 1e7 and sums < 2^25
	for bit := 0; bit < 25; bit++ {
		mod := 1 << (bit + 1)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			b[i] = a[i] % mod
		}
		sort.Ints(b)
		cnt := 0
		for i := 0; i < n; i++ {
			x := b[i]
			// first interval
			l1 := (1 << bit) - x
			r1 := (1 << (bit + 1)) - x
			// second interval
			l2 := (1 << (bit + 1)) + (1 << bit) - x
			r2 := (1 << (bit + 2)) - x

			c := 0
			idx1 := sort.SearchInts(b, l1)
			if idx1 < i+1 {
				idx1 = i + 1
			}
			idx2 := sort.SearchInts(b, r1)
			if idx2 < i+1 {
				idx2 = i + 1
			}
			c += idx2 - idx1
			idx3 := sort.SearchInts(b, l2)
			if idx3 < i+1 {
				idx3 = i + 1
			}
			idx4 := sort.SearchInts(b, r2)
			if idx4 < i+1 {
				idx4 = i + 1
			}
			c += idx4 - idx3
			cnt += c
		}
		if cnt%2 == 1 {
			ans |= 1 << bit
		}
	}

	fmt.Fprintln(out, ans)
}
