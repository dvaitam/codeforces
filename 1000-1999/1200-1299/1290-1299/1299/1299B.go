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
	pts := make([][2]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i][0], &pts[i][1])
	}
	if n%2 == 1 {
		fmt.Fprintln(out, "NO")
		return
	}
	dx := pts[0][0] + pts[n/2][0]
	dy := pts[0][1] + pts[n/2][1]
	for i := 1; i < n/2; i++ {
		if pts[i][0]+pts[i+n/2][0] != dx || pts[i][1]+pts[i+n/2][1] != dy {
			fmt.Fprintln(out, "NO")
			return
		}
	}
	fmt.Fprintln(out, "YES")
}
