package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var V, E int
	if _, err := fmt.Fscan(reader, &V, &E); err != nil {
		return
	}
	x := make([]int, V+1)
	for i := 1; i <= V; i++ {
		fmt.Fscan(reader, &x[i])
	}
	var ans float64
	for i := 0; i < E; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		val := float64(x[u]+x[v]) / float64(w)
		if val > ans {
			ans = val
		}
	}
	fmt.Printf("%.16f\n", ans)
}
