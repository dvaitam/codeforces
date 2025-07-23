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

	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}

	counts := make([]int, 26)
	for i := 0; i < 26 && k > 0; i++ {
		x := 1
		for x*(x-1)/2 <= k {
			x++
		}
		x--
		counts[i] = x
		k -= x * (x - 1) / 2
	}

	used := false
	for i := 0; i < 26; i++ {
		for j := 0; j < counts[i]; j++ {
			writer.WriteByte(byte('a' + i))
			used = true
		}
	}
	if !used {
		writer.WriteByte('a')
	}
}
