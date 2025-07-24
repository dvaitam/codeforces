package main

import (
	"bufio"
	"fmt"
	"os"
)

func digits(x int) int {
	if x == 0 {
		return 1
	}
	d := 0
	for x > 0 {
		d++
		x /= 10
	}
	return d
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var A, B, N int
	fmt.Fscan(in, &A, &B, &N)

	pow10 := make([]int, 11)
	pow10[0] = 1
	for i := 1; i < len(pow10); i++ {
		pow10[i] = pow10[i-1] * 10
	}

	maxK := 1
	for k := 1; k < len(pow10); k++ {
		if pow10[k] <= N {
			maxK = k
		}
	}

	count := 0
	for a := 1; a < A; a++ {
		for b := 1; b < B; b++ {
			m := digits(b)
			den := a*pow10[m] - b
			if den <= 0 {
				continue
			}
			for k := 1; k <= maxK; k++ {
				num := a * b * (pow10[k] - 1)
				if num%den != 0 {
					continue
				}
				n := num / den
				if n >= 1 && n < N && digits(n) == k {
					count++
				}
			}
		}
	}

	fmt.Fprintln(out, count)
}
