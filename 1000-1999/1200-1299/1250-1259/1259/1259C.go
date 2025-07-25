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
		res := solve(s)
		fmt.Fprintln(out, len(res))
		for i, v := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		if len(res) > 0 {
			fmt.Fprintln(out)
		}
	}
}

func solve(s string) []int {
	res := []int{}
	n := len(s)
	i := 0
	for i < n {
		if i+5 <= n && s[i:i+5] == "twone" {
			// remove the middle 'o'
			res = append(res, i+3) // 1-based index
			i += 5
		} else if i+3 <= n && (s[i:i+3] == "one" || s[i:i+3] == "two") {
			// remove the middle letter
			res = append(res, i+2)
			i += 3
		} else {
			i++
		}
	}
	return res
}
