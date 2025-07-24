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
		var s int64
		fmt.Fscan(in, &s)
		var spent int64
		for s >= 10 {
			spent += (s / 10) * 10
			s = s/10 + s%10
		}
		spent += s
		fmt.Fprintln(out, spent)
	}
}
