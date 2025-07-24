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

	var n int
	var s string
	fmt.Fscan(reader, &n)
	fmt.Fscan(reader, &s)

	var res []int
	cnt := 0
	for i := 0; i < n; i++ {
		if s[i] == 'B' {
			cnt++
		} else if cnt > 0 {
			res = append(res, cnt)
			cnt = 0
		}
	}
	if cnt > 0 {
		res = append(res, cnt)
	}

	fmt.Fprintln(writer, len(res))
	if len(res) > 0 {
		for i, x := range res {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, x)
		}
	}
	fmt.Fprintln(writer)
}
