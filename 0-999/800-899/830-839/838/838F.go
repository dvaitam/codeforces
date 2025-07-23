package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var X int64
	if _, err := fmt.Fscan(in, &n, &X); err != nil {
		return
	}
	p := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		fmt.Fscan(in, &p[i])
	}
	var sum int64
	for i := 0; i <= n; i++ {
		sum += int64(i) * p[i]
	}
	c := float64(X) / 1e6
	mu := float64(sum) / 1e6
	ans := mu - c*float64(n)
	if ans < 0 {
		ans = 0
	}
	fmt.Printf("%.10f", ans)
}
