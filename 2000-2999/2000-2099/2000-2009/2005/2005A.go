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
	fmt.Fscan(in, &t)
	vowels := []rune{'a', 'e', 'i', 'o', 'u'}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		q := n / 5
		r := n % 5
		var sb strings.Builder
		for i, ch := range vowels {
			cnt := q
			if i < r {
				cnt++
			}
			for j := 0; j < cnt; j++ {
				sb.WriteRune(ch)
			}
		}
		fmt.Fprintln(out, sb.String())
	}
}
