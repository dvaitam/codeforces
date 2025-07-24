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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	order := make([]int, 0, k)
	onScreen := make(map[int]bool)
	for i := 0; i < n; i++ {
		var id int
		fmt.Fscan(in, &id)
		if onScreen[id] {
			continue
		}
		if len(order) == k {
			last := order[len(order)-1]
			delete(onScreen, last)
			order = order[:len(order)-1]
		}
		// insert at the front
		order = append([]int{id}, order...)
		onScreen[id] = true
	}

	fmt.Fprintln(out, len(order))
	for i, id := range order {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, id)
	}
	if len(order) > 0 {
		out.WriteByte('\n')
	}
}
