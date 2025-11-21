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
		var a1, a2, a4, a5 int
		fmt.Fscan(in, &a1, &a2, &a4, &a5)
		v1 := a1 + a2
		v2 := a4 - a2
		v3 := a5 - a4
		ans := 1
		if v1 == v2 && v2 == v3 {
			ans = 3
		} else if v1 == v2 || v1 == v3 || v2 == v3 {
			ans = 2
		}
		fmt.Fprintln(out, ans)
	}
}
