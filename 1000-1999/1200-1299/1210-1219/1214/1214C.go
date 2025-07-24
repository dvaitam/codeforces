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

	if n%2 != 0 {
		fmt.Fprintln(out, "No")
		return
	}
	bal := 0
	minBal := 0
	for _, ch := range s {
		if ch == '(' {
			bal++
		} else {
			bal--
		}
		if bal < minBal {
			minBal = bal
		}
	}
	if bal == 0 && minBal >= -1 {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}
