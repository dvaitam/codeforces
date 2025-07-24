package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int64
	var k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	ten := int64(1)
	for i := 0; i < k; i++ {
		ten *= 10
	}

	g := gcd(n, ten)
	result := n * (ten / g)
	fmt.Fprintln(writer, result)
}
