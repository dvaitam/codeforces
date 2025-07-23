package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var V int
	if _, err := fmt.Fscan(in, &V); err != nil {
		return
	}
	const n = 250
	for i := 0; i < V; i++ {
		var sum float64
		values := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &values[j])
			sum += float64(values[j])
		}
		mean := sum / n
		var s2 float64
		for j := 0; j < n; j++ {
			diff := float64(values[j]) - mean
			s2 += diff * diff
		}
		variance := s2 / n
		if variance > mean*2 {
			fmt.Fprintln(out, "uniform")
		} else {
			fmt.Fprintln(out, "poisson")
		}
	}
}
