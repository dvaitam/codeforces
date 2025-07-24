package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k, m int
	if _, err := fmt.Fscan(reader, &n, &k, &m); err != nil {
		return
	}
	counts := make(map[int][]int)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		r := x % m
		counts[r] = append(counts[r], x)
		if len(counts[r]) == k {
			fmt.Fprintln(writer, "Yes")
			for j, val := range counts[r] {
				if j > 0 {
					fmt.Fprint(writer, " ")
				}
				fmt.Fprint(writer, val)
			}
			fmt.Fprintln(writer)
			return
		}
	}
	fmt.Fprintln(writer, "No")
}
