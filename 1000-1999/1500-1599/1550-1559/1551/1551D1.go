package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	total := n * m / 2
	k1 := k
	k2 := total - k1
	if n%2 == 1 {
		need := m / 2
		if k1 < need {
			fmt.Fprintln(writer, "NO")
			return
		}
		k1 -= need
		n--
	} else if m%2 == 1 {
		need := n / 2
		if k2 < need {
			fmt.Fprintln(writer, "NO")
			return
		}
		k2 -= need
		m--
	}
	if k1%2 == 0 && k2%2 == 0 {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
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
		solve(reader, writer)
	}
}
