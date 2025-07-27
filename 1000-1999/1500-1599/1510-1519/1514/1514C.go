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
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	result := make([]int, 0)
	prod := 1 % n
	for i := 1; i < n; i++ {
		if gcd(i, n) == 1 {
			result = append(result, i)
			prod = prod * i % n
		}
	}

	if prod != 1 {
		filtered := make([]int, 0, len(result)-1)
		for _, v := range result {
			if v != prod {
				filtered = append(filtered, v)
			}
		}
		result = filtered
	}

	fmt.Fprintln(writer, len(result))
	for i, v := range result {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	if len(result) > 0 {
		fmt.Fprintln(writer)
	}
}
