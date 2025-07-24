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
	var s string
	fmt.Fscan(in, &n)
	fmt.Fscan(in, &s)

	// Trim leading zeros
	i := 0
	for i < len(s) && s[i] == '0' {
		i++
	}
	if i == len(s) {
		fmt.Fprintln(out, 0)
		return
	}
	t := s[i:]
	best := t
	l := len(t)
	maxShift := 20
	if l < maxShift {
		maxShift = l
	}
	for shift := 1; shift <= maxShift; shift++ {
		prefix := t[:l-shift]
		cand := []byte(t)
		for j := shift; j < l; j++ {
			if prefix[j-shift] == '1' {
				cand[j] = '1'
			}
		}
		cs := string(cand)
		if len(cs) > len(best) || cs > best {
			best = cs
		}
	}
	fmt.Fprintln(out, best)
}
