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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		total += a[i]
	}

	match := make([]bool, 101)
	match[0] = true
	copied := 0
	for i := 0; i < n; i++ {
		for x := 1; x <= a[i]; x++ {
			p1 := 100 * x / a[i]
			p2 := 100 * (copied + x) / total
			if p1 == p2 {
				match[p1] = true
			}
		}
		copied += a[i]
	}

	first := true
	for i := 0; i <= 100; i++ {
		if match[i] {
			if !first {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, i)
			first = false
		}
	}
	writer.WriteByte('\n')
}
