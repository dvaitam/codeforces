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
		balance := 0
		ok := false
		for i := 0; i < len(s); i++ {
			if s[i] == '(' {
				balance++
			} else {
				balance--
			}
			if balance == 0 && i != len(s)-1 {
				ok = true
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
