package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k [3]int
	fmt.Fscan(in, &k[0], &k[1], &k[2])
	sort.Ints(k[:])

	if k[0] == 1 || (k[0] == 2 && k[1] == 2) ||
		(k[0] == 2 && k[1] == 4 && k[2] == 4) ||
		(k[0] == 3 && k[1] == 3 && k[2] == 3) {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
