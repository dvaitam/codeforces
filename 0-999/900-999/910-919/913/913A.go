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

	var n, m int64
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	if n >= 31 {
		fmt.Fprintln(out, m)
	} else {
		mod := int64(1) << uint(n)
		fmt.Fprintln(out, m%mod)
	}
}
