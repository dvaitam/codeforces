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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var sb strings.Builder
		ones := k
		zeros := n - k
		if ones > 0 {
			sb.WriteString(strings.Repeat("1", ones))
			sb.WriteString(strings.Repeat("0", zeros))
		} else {
			sb.WriteString(strings.Repeat("0", n))
		}
		fmt.Fprintln(out, sb.String())
	}
}
