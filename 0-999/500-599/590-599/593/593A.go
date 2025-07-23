package main

import (
	"bufio"
	"fmt"
	"math/bits"
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

	type word struct {
		mask uint32
		len  int
	}
	words := make([]word, 0, n)

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		var m uint32
		for _, ch := range s {
			m |= 1 << (ch - 'a')
		}
		if bits.OnesCount32(m) <= 2 {
			words = append(words, word{m, len(s)})
		}
	}

	maxLen := 0
	for i := 0; i < 26; i++ {
		for j := i; j < 26; j++ {
			pairMask := uint32((1 << i) | (1 << j))
			sum := 0
			for _, w := range words {
				if w.mask & ^pairMask == 0 {
					sum += w.len
				}
			}
			if sum > maxLen {
				maxLen = sum
			}
		}
	}

	fmt.Fprintln(writer, maxLen)
}
