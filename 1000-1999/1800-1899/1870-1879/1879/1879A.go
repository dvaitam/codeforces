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
		s := make([]int64, n)
		e := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &s[i], &e[i])
		}
		s1 := s[0]
		e1 := e[0]
		var maxS int64 = -1
		for i := 1; i < n; i++ {
			if e[i] >= e1 {
				if s[i] > maxS {
					maxS = s[i]
				}
			}
		}
		if s1 > maxS {
			fmt.Fprintln(writer, s1)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
