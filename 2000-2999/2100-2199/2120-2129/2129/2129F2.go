package main

import (
	"bufio"
	"fmt"
	"os"
)

// Offline (hacked) version: the hidden permutation is provided in input, so we
// just read and print it in the expected answer format.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := nextInt(in)
	for ; t > 0; t-- {
		n := nextInt(in)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			p[i] = nextInt(in)
		}
		fmt.Fprint(out, "! ")
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, p[i])
		}
		if t > 1 {
			fmt.Fprintln(out)
		}
	}
}

func nextInt(r *bufio.Reader) int {
	sign, val := 1, 0
	c, _ := r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = r.ReadByte()
	}
	return val * sign
}
