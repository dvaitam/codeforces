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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		minIdx := 0
		for i := 1; i < len(s); i++ {
			if s[i] < s[minIdx] {
				minIdx = i
			}
		}
		a := string(s[minIdx])
		b := s[:minIdx] + s[minIdx+1:]
		fmt.Fprintln(out, a, b)
	}
}
