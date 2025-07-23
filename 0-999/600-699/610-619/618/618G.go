package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	var p int64
	if _, err := fmt.Fscan(in, &n, &p); err != nil {
		return
	}
	prob1 := float64(p) / 1e9
	expected := float64(n * (n + 1) / 2)
	expected += 1 - prob1
	fmt.Printf("%.6f\n", expected)
}
