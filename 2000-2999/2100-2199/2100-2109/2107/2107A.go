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
		a := make([]int, n)
		same := true
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if i > 0 && a[i] != a[0] {
				same = false
			}
		}
		if same {
			fmt.Fprintln(out, "No")
			continue
		}
		fmt.Fprintln(out, "Yes")
		for i := 0; i < n; i++ {
			if a[i] == a[0] {
				fmt.Fprint(out, "2 ")
			} else {
				fmt.Fprint(out, "1 ")
			}
		}
		fmt.Fprintln(out)
	}
}
