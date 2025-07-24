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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([][]int, n)
	for i := 0; i < n; i++ {
		p[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &p[i][j])
		}
	}

	// TODO: implement probability calculation

	fmt.Fprintln(out, 0)
}
