package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)

	best := ""
	for i := 1; i <= n; i++ {
		pref := s[:i]
		var b strings.Builder
		for b.Len() < k {
			b.WriteString(pref)
		}
		cand := b.String()[:k]
		if best == "" || cand < best {
			best = cand
		}
	}
	fmt.Println(best)
}
