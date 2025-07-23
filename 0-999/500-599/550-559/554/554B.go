package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	counts := make(map[string]int, n)
	for i := 0; i < n; i++ {
		var row string
		fmt.Fscan(in, &row)
		counts[row]++
	}
	max := 0
	for _, c := range counts {
		if c > max {
			max = c
		}
	}
	fmt.Println(max)
}
