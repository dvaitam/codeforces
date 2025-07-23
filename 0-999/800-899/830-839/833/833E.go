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
	var C int64
	if _, err := fmt.Fscan(reader, &n, &C); err != nil {
		return
	}
	type cloud struct {
		l, r int64
		c    int64
	}
	clouds := make([]cloud, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &clouds[i].l, &clouds[i].r, &clouds[i].c)
	}

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	for i := 0; i < m; i++ {
		var k int64
		fmt.Fscan(reader, &k)
		// TODO: implement algorithm described in problemE.txt
		// Placeholder implementation prints 0 for each query
		fmt.Fprintln(writer, 0)
	}
}
