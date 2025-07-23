package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, t int
	if _, err := fmt.Fscan(in, &n, &t); err != nil {
		return
	}
	g := make([][]float64, n+2)
	for i := range g {
		g[i] = make([]float64, n+2)
	}
	g[1][1] = float64(t)
	for i := 1; i <= n; i++ {
		for j := 1; j <= i; j++ {
			if g[i][j] > 1.0 {
				overflow := g[i][j] - 1.0
				g[i][j] = 1.0
				g[i+1][j] += overflow / 2
				g[i+1][j+1] += overflow / 2
			}
		}
	}
	ans := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= i; j++ {
			if g[i][j] >= 1.0-1e-9 {
				ans++
			}
		}
	}
	fmt.Println(ans)
}
