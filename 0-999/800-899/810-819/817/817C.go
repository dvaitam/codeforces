package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumDigits(x int64) int64 {
	var s int64
	for x > 0 {
		s += x % 10
		x /= 10
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, s int64
	if _, err := fmt.Fscan(in, &n, &s); err != nil {
		return
	}

	low, high := int64(1), n
	for low <= high {
		mid := (low + high) / 2
		if mid-sumDigits(mid) >= s {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	if low > n {
		fmt.Fprintln(out, 0)
	} else {
		fmt.Fprintln(out, n-low+1)
	}
}
