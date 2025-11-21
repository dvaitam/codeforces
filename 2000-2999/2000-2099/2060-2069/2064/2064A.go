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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		var pos1 []int
		var pos0 []int
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				pos1 = append(pos1, i)
			} else {
				pos0 = append(pos0, i)
			}
		}
		if len(pos1) == 0 || len(pos0) == 0 || pos1[len(pos1)-1] < pos0[0] {
			fmt.Fprintln(out, 0)
			continue
		}
		fmt.Fprintln(out, 2)
	}
}
