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

	digits := "0123456789"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower := "abcdefghijklmnopqrstuvwxyz"

	for ; t > 0; t-- {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)

		var sb strings.Builder
		for i := 0; i < a; i++ {
			sb.WriteByte(digits[i%len(digits)])
		}
		for i := 0; i < b; i++ {
			sb.WriteByte(upper[i%len(upper)])
		}
		for i := 0; i < c; i++ {
			sb.WriteByte(lower[i%len(lower)])
		}

		fmt.Fprintln(out, sb.String())
	}
}
