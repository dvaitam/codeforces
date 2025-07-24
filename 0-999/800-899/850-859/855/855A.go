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

	seen := make(map[string]struct{})
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		if _, ok := seen[s]; ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
			seen[s] = struct{}{}
		}
	}
}
