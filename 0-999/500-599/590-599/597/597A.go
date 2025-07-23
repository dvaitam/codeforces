package main

import (
	"bufio"
	"fmt"
	"os"
)

// floorDiv computes floor(a/b) for positive b
func floorDiv(a, b int64) int64 {
	if a >= 0 {
		return a / b
	}
	// adjust for negative a to perform floor division
	return -((-a + b - 1) / b)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var k, a, b int64
	if _, err := fmt.Fscan(in, &k, &a, &b); err != nil {
		return
	}
	result := floorDiv(b, k) - floorDiv(a-1, k)
	fmt.Println(result)
}
