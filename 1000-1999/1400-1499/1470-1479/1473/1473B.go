package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		l := lcm(len(s), len(t))
		rs := strings.Repeat(s, l/len(s))
		rt := strings.Repeat(t, l/len(t))
		if rs == rt {
			fmt.Fprintln(writer, rs)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
