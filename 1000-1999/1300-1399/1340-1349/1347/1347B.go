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
		var a1, b1, a2, b2 int
		fmt.Fscan(in, &a1, &b1)
		fmt.Fscan(in, &a2, &b2)
		ok := false
		if a1 == a2 && b1+b2 == a1 {
			ok = true
		}
		if a1 == b2 && b1+a2 == a1 {
			ok = true
		}
		if b1 == a2 && a1+b2 == b1 {
			ok = true
		}
		if b1 == b2 && a1+a2 == b1 {
			ok = true
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
