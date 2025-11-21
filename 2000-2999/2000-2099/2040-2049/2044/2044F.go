package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxVal = 200000
const offset = MaxVal
const size = 2*MaxVal + 1

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}

	a := make([]int64, n)
	var sumA int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		sumA += a[i]
	}

	b := make([]int64, m)
	var sumB int64
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &b[i])
		sumB += b[i]
	}

	hasA := make([]bool, size)
	hasB := make([]bool, size)

	for i := 0; i < n; i++ {
		diff := sumA - a[i]
		if absInt64(diff) <= MaxVal {
			hasA[int(diff)+offset] = true
		}
	}
	for i := 0; i < m; i++ {
		diff := sumB - b[i]
		if absInt64(diff) <= MaxVal {
			hasB[int(diff)+offset] = true
		}
	}

	check := func(x int) bool {
		if x == 0 {
			return false
		}
		idx := x + offset
		if idx < 0 || idx >= size {
			return false
		}
		return hasA[idx]
	}
	checkB := func(x int) bool {
		if x == 0 {
			return false
		}
		idx := x + offset
		if idx < 0 || idx >= size {
			return false
		}
		return hasB[idx]
	}

	for ; q > 0; q-- {
		var x int
		fmt.Fscan(in, &x)
		ans := false
		absx := absInt(x)
		for d := 1; d*d <= absx && !ans; d++ {
			if absx%d != 0 {
				continue
			}
			divisors := []int{d}
			other := absx / d
			if other != d {
				divisors = append(divisors, other)
			}
			for _, val := range divisors {
				for _, s := range []int{val, -val} {
					if s == 0 {
						continue
					}
					if !check(s) {
						continue
					}
					if x%s != 0 {
						continue
					}
					t := x / s
					if checkB(t) {
						ans = true
						break
					}
				}
				if ans {
					break
				}
			}
		}
		if ans {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
