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
		var n int
		fmt.Fscan(in, &n)
		uniq := int64(-1)
		ok := true
		for i := 0; i < n; i++ {
			var val int64
			fmt.Fscan(in, &val)
			if val == -1 {
				continue
			}
			if val == 0 {
				ok = false
			}
			if uniq == -1 {
				uniq = val
			} else if val != uniq {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
