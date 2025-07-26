package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		var sb strings.Builder
		for i := 0; i < n; i++ {
			for j := 0; j < k; j++ {
				sb.WriteByte(byte('a' + j))
			}
		}
		fmt.Fprintln(out, sb.String())
	}
}
