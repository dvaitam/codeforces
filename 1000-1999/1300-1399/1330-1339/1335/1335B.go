package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, a, b int
		fmt.Fscan(reader, &n, &a, &b)
		s := make([]byte, 0, n)
		// append first b distinct characters
		for i := 0; i < b; i++ {
			s = append(s, byte('a'+i))
		}
		idx := 0
		for len(s) < n {
			idx %= a
			s = append(s, s[idx])
			idx++
		}
		writer.Write(s)
		writer.WriteByte('\n')
	}
}
