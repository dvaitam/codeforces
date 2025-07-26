package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < n; i++ {
		out.WriteByte('1')
	}
	out.WriteByte('\n')
}
