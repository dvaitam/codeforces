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
		var n int
		fmt.Fscan(reader, &n)
		s1 := make([]byte, n)
		s2 := make([]byte, n)
		fmt.Fscan(reader, &s1)
		fmt.Fscan(reader, &s2)
		ok := true
		for i := 0; i < n; i++ {
			a := s1[i] == 'R'
			b := s2[i] == 'R'
			if a != b {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
