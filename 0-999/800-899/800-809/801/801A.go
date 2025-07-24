package main

import (
	"bufio"
	"fmt"
	"os"
)

func countVK(b []byte) int {
	count := 0
	for i := 0; i+1 < len(b); i++ {
		if b[i] == 'V' && b[i+1] == 'K' {
			count++
		}
	}
	return count
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	b := []byte(s)
	best := countVK(b)

	for i := 0; i < len(b); i++ {
		orig := b[i]
		if b[i] == 'V' {
			b[i] = 'K'
		} else {
			b[i] = 'V'
		}
		if c := countVK(b); c > best {
			best = c
		}
		b[i] = orig
	}

	fmt.Fprintln(out, best)
}
