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
		var s string
		fmt.Fscan(in, &s)
		cntA := 0
		for i := 0; i < len(s); i++ {
			if s[i] == 'A' {
				cntA++
			}
		}
		if cntA > len(s)-cntA {
			fmt.Fprintln(out, "A")
		} else {
			fmt.Fprintln(out, "B")
		}
	}
}
