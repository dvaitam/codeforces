package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		countT := 0
		var others bytes.Buffer
		for i := 0; i < len(s); i++ {
			if s[i] == 'T' {
				countT++
			} else {
				others.WriteByte(s[i])
			}
		}
		result := bytes.Repeat([]byte{'T'}, countT)
		result = append(result, others.Bytes()...)
		fmt.Fprintln(out, string(result))
	}
}
