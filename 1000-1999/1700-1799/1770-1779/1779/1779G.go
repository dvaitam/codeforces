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
		var n int
		var s1, s2, s3 string
		fmt.Fscan(in, &n, &s1, &s2, &s3)
		fmt.Fprintln(out, 0)
	}
}
