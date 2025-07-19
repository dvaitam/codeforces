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
		var k int64
		fmt.Fscan(in, &k)
		if k&1 == 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		// decode bits
		var tag [62]bool
		for i := 0; i < 62; i++ {
			if (k>>i)&1 == 1 {
				tag[i] = true
			}
		}
		// build reversed checkpoint flags
		ans := make([]bool, 1)
		ptr := 0
		// initial dummy
		ans[0] = true
		// bit 1 special
		if tag[1] {
			ptr++
			if ptr >= len(ans) {
				ans = append(ans, true)
			} else {
				ans[ptr] = true
			}
		}
		// bits >=2
		for i := 2; i < 62; i++ {
			if tag[i] {
				ptr++
				if ptr >= len(ans) {
					ans = append(ans, true)
				} else {
					ans[ptr] = true
				}
				ptr += i - 1
				// ensure capacity
				for ptr >= len(ans) {
					ans = append(ans, false)
				}
				ans[ptr] = true
			}
		}
		// output
		fmt.Fprintln(out, ptr)
		for i := ptr; i >= 1; i-- {
			if ans[i] {
				out.WriteString("1 ")
			} else {
				out.WriteString("0 ")
			}
		}
		out.WriteByte('\n')
	}
}
