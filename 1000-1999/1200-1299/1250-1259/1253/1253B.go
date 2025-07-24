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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	events := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &events[i])
	}

	inside := make(map[int]bool)
	seen := make(map[int]bool)
	res := make([]int, 0)
	last := 0

	for i, x := range events {
		if x > 0 {
			if seen[x] {
				fmt.Fprintln(writer, -1)
				return
			}
			seen[x] = true
			inside[x] = true
		} else {
			id := -x
			if !inside[id] {
				fmt.Fprintln(writer, -1)
				return
			}
			delete(inside, id)
		}
		if len(inside) == 0 {
			res = append(res, i-last+1)
			last = i + 1
			seen = make(map[int]bool)
		}
	}

	if len(inside) != 0 || last != n {
		fmt.Fprintln(writer, -1)
		return
	}

	fmt.Fprintln(writer, len(res))
	for i, v := range res {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}
