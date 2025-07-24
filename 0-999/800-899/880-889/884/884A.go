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

	var n int
	var t int64
	if _, err := fmt.Fscan(in, &n, &t); err != nil {
		return
	}

	for i := 1; i <= n; i++ {
		var a int64
		fmt.Fscan(in, &a)
		t -= 86400 - a
		if t <= 0 {
			fmt.Fprintln(out, i)
			return
		}
	}
}
