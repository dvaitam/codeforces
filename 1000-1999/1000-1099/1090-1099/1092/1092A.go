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

	var t, n, k int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		fmt.Fscan(reader, &n, &k)
		p := n / k
		base := byte('a')
		// print k groups of p letters
		for i := 0; i < k; i++ {
			ch := base + byte(i)
			for j := 0; j < p; j++ {
				writer.WriteByte(ch)
			}
		}
		// print remaining letters as last character
		rem := n - k*p
		if rem > 0 {
			last := base + byte(k-1)
			for i := 0; i < rem; i++ {
				writer.WriteByte(last)
			}
		}
		writer.WriteByte('\n')
	}
}
