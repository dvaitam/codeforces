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
		p := make([]int, n)
		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			pos[p[i]] = i
		}
		ok := true
		for i := 1; i < n; i++ {
			if pos[i+1] > pos[i] && pos[i+1] != pos[i]+1 {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
