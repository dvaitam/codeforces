package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	r := make([]int, n)
	s := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &r[i])
		s += r[i]
	}
	if m == 1 {
		// Only one color, trivial
		fmt.Println("1.000000000")
		return
	}
	maxJ := n * m
	prev := make([]float64, maxJ+1)
	curr := make([]float64, maxJ+1)
	prev[0] = 1.0
	for i := 1; i <= n; i++ {
		sum := prev[0]
		for j := 1; j <= i*m; j++ {
			if j-m-1 >= 0 {
				sum -= prev[j-m-1]
			}
			val := sum
			if j >= r[i-1] {
				val -= prev[j-r[i-1]]
			}
			curr[j] = val / float64(m-1)
			sum += prev[j]
		}
		// move curr into prev and reset curr
		for j := 0; j <= i*m; j++ {
			prev[j] = curr[j]
			curr[j] = 0
		}
	}
	ans := 0.0
	for j := 0; j < s; j++ {
		ans += prev[j]
	}
	ans = ans*float64(m-1) + 1.0
	fmt.Printf("%.9f\n", ans)
}
