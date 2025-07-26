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
	target := "abc"
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		diff := make([]int, 0, 2)
		for i := 0; i < 3; i++ {
			if s[i] != target[i] {
				diff = append(diff, i)
			}
		}
		ok := false
		if len(diff) == 0 {
			ok = true
		} else if len(diff) == 2 {
			i, j := diff[0], diff[1]
			if s[i] == target[j] && s[j] == target[i] {
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
