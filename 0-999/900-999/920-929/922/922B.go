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

	var count int64
	for a := 1; a <= n; a++ {
		for b := a; b <= n; b++ {
			c := a ^ b
			if c < b || c > n {
				continue
			}
			if a+b <= c {
				continue
			}
			count++
		}
	}

	fmt.Fprintln(out, count)
}
