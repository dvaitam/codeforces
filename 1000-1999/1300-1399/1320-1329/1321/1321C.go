package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	bytes := []byte(s)
	removed := 0
	for ch := byte('z'); ch > 'a'; ch-- {
		for {
			idx := -1
			for i := 0; i < len(bytes); i++ {
				if bytes[i] == ch {
					if (i > 0 && bytes[i-1] == ch-1) || (i+1 < len(bytes) && bytes[i+1] == ch-1) {
						idx = i
						break
					}
				}
			}
			if idx == -1 {
				break
			}
			bytes = append(bytes[:idx], bytes[idx+1:]...)
			removed++
		}
	}

	fmt.Fprintln(out, removed)
}
