package main

import (
	"bufio"
	"fmt"
	"os"
)

type block struct {
	l, r int
}

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

		var s string
		fmt.Fscan(reader, &s)
		bytes := []byte(s)

		var blocks []block
		possible := true

		for i := 0; i < n && possible; {
			if bytes[i] == '0' {
				j := i
				for j < n && bytes[j] == '0' {
					j++
				}
				if j-i == 1 {
					possible = false
					break
				}
				blocks = append(blocks, block{i + 1, j})
				i = j
			} else {
				i++
			}
		}

		if !possible {
			fmt.Fprintln(writer, "NO")
			continue
		}

		p := make([]int, n)
		for i := 0; i < n; i++ {
			if bytes[i] == '1' {
				p[i] = i + 1
			}
		}

		for _, b := range blocks {
			for pos := b.l; pos < b.r; pos++ {
				p[pos-1] = pos + 1
			}
			p[b.r-1] = b.l
		}

		fmt.Fprintln(writer, "YES")
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, p[i])
		}
		fmt.Fprintln(writer)
	}
}
