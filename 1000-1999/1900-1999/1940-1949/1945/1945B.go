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
		var a, b, m int64
		fmt.Fscan(in, &a, &b, &m)
		ans := m/a + m/b + 2
		fmt.Fprintln(out, ans)
	}
}
