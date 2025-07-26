package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func check(n int64) bool {
	if n < 7 { // minimum possible is 1+2+4=7
		return false
	}
	for h := 2; h <= 60; h++ {
		base := int64(math.Pow(float64(n), 1.0/float64(h)))
		if base < 2 {
			base = 2
		}
		for delta := int64(-1); delta <= 1; delta++ {
			k := base + delta
			if k < 2 {
				continue
			}
			sum := int64(1)
			cur := int64(1)
			overflow := false
			for i := 1; i <= h; i++ {
				if cur > n/k {
					overflow = true
					break
				}
				cur *= k
				sum += cur
				if sum > n {
					overflow = true
					break
				}
			}
			if !overflow && sum == n {
				return true
			}
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		if check(n) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
