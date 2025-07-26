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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		k := n / 2
		ops := 0
		pred := make([]byte, 0, k)
		for i := 0; i < k; i++ {
			a := s[2*i]
			b := s[2*i+1]
			if a != b {
				ops++
			} else {
				pred = append(pred, a)
			}
		}
		seg := 1
		if len(pred) > 0 {
			seg = 1
			for i := 1; i < len(pred); i++ {
				if pred[i] != pred[i-1] {
					seg++
				}
			}
		}
		fmt.Fprintln(out, ops, seg)
	}
}
