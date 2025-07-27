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
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		m := a
		if b > m {
			m = b
		}
		if c > m {
			m = c
		}
		cnt := 0
		if a == m {
			cnt++
		}
		if b == m {
			cnt++
		}
		if c == m {
			cnt++
		}
		ansA := m + 1 - a
		ansB := m + 1 - b
		ansC := m + 1 - c
		if cnt == 1 {
			if a == m {
				ansA = 0
			}
			if b == m {
				ansB = 0
			}
			if c == m {
				ansC = 0
			}
		}
		fmt.Fprintln(writer, ansA, ansB, ansC)
	}
}
