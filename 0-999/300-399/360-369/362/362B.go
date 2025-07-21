package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	rdr := bufio.NewReader(os.Stdin)
	wr := bufio.NewWriter(os.Stdout)
	defer wr.Flush()

	var n, m int
	fmt.Fscan(rdr, &n, &m)

	dirty := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(rdr, &dirty[i])
	}

	sort.Ints(dirty)

	if m > 0 && (dirty[0] == 1 || dirty[m-1] == n) {
		fmt.Fprintln(wr, "NO")
		return
	}

	for i := 0; i+2 < m; i++ {
		if dirty[i+2]-dirty[i] == 2 {
			fmt.Fprintln(wr, "NO")
			return
		}
	}

	fmt.Fprintln(wr, "YES")
}
