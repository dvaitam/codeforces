package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const maxN = 500000

var (
	special    []int
	impossible [maxN + 1]bool
)

func isPerfectSquare(x int64) bool {
	if x < 0 {
		return false
	}
	r := int64(math.Sqrt(float64(x)))
	return r*r == x || (r+1)*(r+1) == x
}

func precompute() {
	for i := 1; i <= maxN; i++ {
		sum := int64(i) * int64(i+1) / 2
		if isPerfectSquare(sum) {
			impossible[i] = true
			special = append(special, i)
		}
	}
}

func main() {
	precompute()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if impossible[n] {
			fmt.Fprintln(out, -1)
			continue
		}

		perm := make([]int, n)
		for i := range perm {
			perm[i] = i + 1
		}

		for _, idx := range special {
			if idx >= n {
				break
			}
			perm[idx-1], perm[idx] = perm[idx], perm[idx-1]
		}

		for i, v := range perm {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
