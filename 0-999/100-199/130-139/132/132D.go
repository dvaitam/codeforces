package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return
	}
	str = strings.TrimSpace(str)
	n := len(str)
	// prepare buffer with two leading zeros
	s := make([]byte, n+2)
	s[0], s[1] = '0', '0'
	for i := 0; i < n; i++ {
		s[i+2] = str[i]
	}
	type op struct {
		pow int
		op  byte
	}
	ops := make([]op, 0)
	// process bits from LSB (index n+1) down to 1
	for i := n + 1; i >= 1; i-- {
		if s[i] == '1' {
			if s[i-1] == '0' {
				ops = append(ops, op{pow: n + 1 - i, op: '+'})
			} else {
				// handle carry: flip run of ones to zero and set next bit
				j := i
				for j >= 0 && s[j] == '1' {
					s[j] = '0'
					j--
				}
				s[j] = '1'
				ops = append(ops, op{pow: n + 1 - i, op: '-'})
			}
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, len(ops))
	for _, o := range ops {
		fmt.Fprintf(writer, "%c2^%d\n", o.op, o.pow)
	}
}
