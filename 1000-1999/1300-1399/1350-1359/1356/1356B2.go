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
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if len(s) < n {
		tmp := make([]byte, n)
		copy(tmp, s)
		for i := len(s); i < n; i++ {
			tmp[i] = '0'
		}
		s = string(tmp)
	}
	bits := []byte(s)
	carry := byte(1)
	for i := 0; i < n; i++ {
		if carry == 0 {
			break
		}
		if bits[i] == '0' {
			bits[i] = '1'
			carry = 0
		} else {
			bits[i] = '0'
		}
	}
	if carry == 1 {
		for i := 0; i < n; i++ {
			bits[i] = '1'
		}
	}
	fmt.Fprintln(out, string(bits))
}
