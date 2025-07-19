package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	res := make([]byte, n)
	p, q := true, true
	for i := 0; i < n && i < len(s); i++ {
		if s[i] == '(' {
			if p {
				res[i] = '0'
			} else {
				res[i] = '1'
			}
			p = !p
		} else {
			if q {
				res[i] = '0'
			} else {
				res[i] = '1'
			}
			q = !q
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	writer.Write(res)
	writer.WriteByte('\n')
	writer.Flush()
}
