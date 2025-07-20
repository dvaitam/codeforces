package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1<<20)  // 1 MiB buffer for big input
	writer := bufio.NewWriterSize(os.Stdout, 1<<20) // 1 MiB buffer for fast output
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		var s, t string
		fmt.Fscan(reader, &n, &s, &t)

		f1 := strings.IndexByte(s, '1')
		if f1 == -1 {
			f1 = n // no '1' found → behave like pointer-to-end in C++
		}
		f2 := strings.IndexByte(t, '1')
		if f2 == -1 {
			f2 = n
		}

		if f1 <= f2 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

