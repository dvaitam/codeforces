package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var p int64
	if _, err := fmt.Fscan(in, &n, &p); err != nil {
		return
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}
	// TODO: implement dynamic programming solution for problem D
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 1; i <= n-1; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, 0)
	}
	out.WriteByte('\n')
}
