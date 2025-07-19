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

	var tt int
	fmt.Fscan(reader, &tt)
	for t := 0; t < tt; t++ {
		solve(reader, writer)
		writer.WriteByte('\n')
	}
}

// solve processes one test case
func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	a := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		a[i] = []byte(s)
	}
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if a[i][j] == 'X' {
				cnt++
			}
		}
	}
	// try three residue classes
	for k := 0; k < 3; k++ {
		c := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if (i+j)%3 == k && a[i][j] == 'X' {
					c++
				}
			}
		}
		if c <= cnt/3 {
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					if (i+j)%3 == k && a[i][j] == 'X' {
						writer.WriteByte(byte('O'))
					} else {
						writer.WriteByte(a[i][j])
					}
				}
				writer.WriteByte('\n')
			}
			return
		}
	}
}
