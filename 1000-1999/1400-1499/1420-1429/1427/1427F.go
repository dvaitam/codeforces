package main

import (
	"bufio"
	"fmt"
	"os"
)

type elem struct {
	id, col int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	s := n * 6
	// 1-based indexing: allocate s+1
	a := make([]elem, s+1)
	for i := 1; i <= s; i++ {
		a[i].id = i
	}
	// read s/2 marked elements
	for i := 0; i < s/2; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x >= 1 && x <= s {
			a[x].col = 1
		}
	}
	// perform 2n iterations
	for T := 2*n - 1; T >= 0; T-- {
		// current number of elements is s
		// find starting index
		i := 1
		for i <= s-3 && ((i%3) == 0 || a[i].col == a[i+1].col) {
			i++
		}
		// i now is candidate start of a block of 3 (*), move to third pos
		i += 3
		// find triple ending at i
		for i <= s && (((a[i].col+T)&1) != 0 || (a[i-1].col^a[i].col) != 0 || (a[i-2].col^a[i].col) != 0) {
			i++
		}
		if i > s {
			i = 3
			for ((a[i].col+T)&1) != 0 || (a[i-1].col^a[i].col) != 0 || (a[i-2].col^a[i].col) != 0 {
				i++
			}
		}
		// output ids of the chosen triple
		fmt.Fprintf(writer, "%d %d %d\n", a[i-2].id, a[i-1].id, a[i].id)
		// remove these three elements: shift left
		for j := i + 1; j <= s; j++ {
			a[j-3] = a[j]
		}
		s -= 3
	}
}
