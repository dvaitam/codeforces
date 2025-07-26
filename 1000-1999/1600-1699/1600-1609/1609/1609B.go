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

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	var s string
	fmt.Fscan(in, &s)
	b := []byte(s)

	count := 0
	if n >= 3 {
		for i := 0; i+2 < n; i++ {
			if b[i] == 'a' && b[i+1] == 'b' && b[i+2] == 'c' {
				count++
			}
		}
	}

	for ; q > 0; q-- {
		var pos int
		var ch string
		fmt.Fscan(in, &pos, &ch)
		pos--
		// remove occurrences affected by the position before change
		for i := pos - 2; i <= pos; i++ {
			if i >= 0 && i+2 < n {
				if b[i] == 'a' && b[i+1] == 'b' && b[i+2] == 'c' {
					count--
				}
			}
		}
		b[pos] = ch[0]
		// add occurrences after change
		for i := pos - 2; i <= pos; i++ {
			if i >= 0 && i+2 < n {
				if b[i] == 'a' && b[i+1] == 'b' && b[i+2] == 'c' {
					count++
				}
			}
		}
		fmt.Fprintln(out, count)
	}
}
