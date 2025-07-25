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
		total0, total1 := 0, 0
		for i := 0; i < len(s); i++ {
			if s[i] == '0' {
				total0++
			} else {
				total1++
			}
		}
		pref0, pref1, best := 0, 0, 0
		for i := 0; i < len(s); i++ {
			if s[i] == '0' {
				pref0++
			} else {
				pref1++
			}
			if pref0 <= total1 && pref1 <= total0 {
				best = i + 1
			}
		}
		fmt.Fprintln(out, len(s)-best)
	}
}
