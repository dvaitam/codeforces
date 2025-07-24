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

	var n, m int
	fmt.Fscan(in, &n, &m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var x int
			fmt.Fscan(in, &x)
		}
	}

	line := make([]byte, m)
	for j := 0; j < m; j++ {
		line[j] = '0'
	}
	for i := 0; i < n; i++ {
		fmt.Fprintln(out, string(line))
	}
}
