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
		if len(s) == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		blocks := 1
		for i := 1; i < len(s); i++ {
			if s[i] != s[i-1] {
				blocks++
			}
		}
		fmt.Fprintln(out, blocks)
	}
}
