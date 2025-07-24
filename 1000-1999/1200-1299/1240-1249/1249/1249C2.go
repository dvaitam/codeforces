package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(n int64) int64 {
	digits := make([]int, 0, 60)
	for n > 0 {
		digits = append(digits, int(n%3))
		n /= 3
	}
	digits = append(digits, 0)
	for i := 0; i < len(digits); i++ {
		if digits[i] >= 2 {
			for j := 0; j <= i; j++ {
				digits[j] = 0
			}
			if i+1 == len(digits) {
				digits = append(digits, 1)
			} else {
				digits[i+1]++
			}
		}
	}
	var res int64
	pow := int64(1)
	for _, d := range digits {
		if d == 1 {
			res += pow
		}
		pow *= 3
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var n int64
		fmt.Fscan(in, &n)
		fmt.Fprintln(out, solve(n))
	}
}
