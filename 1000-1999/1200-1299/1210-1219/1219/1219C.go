package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// addOne returns the decimal string representing s + 1.
func addOne(s string) string {
	b := []byte(s)
	i := len(b) - 1
	carry := byte(1)
	for i >= 0 && carry > 0 {
		sum := b[i] - '0' + carry
		b[i] = sum%10 + '0'
		carry = sum / 10
		i--
	}
	if carry > 0 {
		b = append([]byte{'1'}, b...)
	}
	return string(b)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var L int
	if _, err := fmt.Fscan(in, &L); err != nil {
		return
	}
	var A string
	if _, err := fmt.Fscan(in, &A); err != nil {
		return
	}

	n := len(A)
	base := "1" + strings.Repeat("0", L-1)

	if n < L {
		fmt.Fprint(out, base)
		return
	}

	if n%L != 0 {
		k := (n + L - 1) / L
		fmt.Fprint(out, strings.Repeat(base, k))
		return
	}

	k := n / L
	prefix := A[:L]
	candidate := strings.Repeat(prefix, k)
	if candidate > A {
		fmt.Fprint(out, candidate)
		return
	}

	inc := addOne(prefix)
	if len(inc) > L {
		fmt.Fprint(out, strings.Repeat(base, k+1))
		return
	}
	fmt.Fprint(out, strings.Repeat(inc, k))
}
