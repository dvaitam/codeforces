package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func buildDivisors(limit int) [][]int {
	divs := make([][]int, limit+1)
	for d := 1; d <= limit; d++ {
		for multiple := d; multiple <= limit; multiple += d {
			divs[multiple] = append(divs[multiple], d)
		}
	}
	return divs
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		limit := n
		divs := buildDivisors(limit)
		cnt := make([]int, limit+1)
		gcdVal := make([]int, limit+1)
		good := make([]bool, limit+1)

		res := make([]int, n)
		globalG := 0
		best := 0

		for i := 0; i < n; i++ {
			x := arr[i]
			oldG := globalG
			if globalG == 0 {
				globalG = x
			} else {
				globalG = gcd(globalG, x)
			}

			if oldG != 0 && globalG != oldG {
				for _, d := range divs[oldG] {
					if d <= 1 {
						continue
					}
					if globalG%d == 0 {
						continue
					}
					if good[d] && cnt[d] > best {
						best = cnt[d]
					}
				}
			}

			for _, d := range divs[x] {
				if d <= 1 {
					continue
				}
				cnt[d]++
				if gcdVal[d] == 0 {
					gcdVal[d] = x
				} else {
					gcdVal[d] = gcd(gcdVal[d], x)
				}
				if !good[d] && gcdVal[d] == d {
					good[d] = true
				}
				if good[d] && globalG%d != 0 && cnt[d] > best {
					best = cnt[d]
				}
			}

			res[i] = best
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[i])
		}
		fmt.Fprintln(out)
	}
}
