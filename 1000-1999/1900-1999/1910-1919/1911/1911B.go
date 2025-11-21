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
	var s string
	if _, err := fmt.Fscan(in, &n, &s); err != nil {
		return
	}

	res := make([]byte, 0, n)
	for i, step := 0, 1; i < n; i += step {
		res = append(res, s[i])
		step++
	}

	fmt.Fprintln(out, string(res))
}
