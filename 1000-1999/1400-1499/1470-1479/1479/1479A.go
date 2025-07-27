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

	values := map[int]int{}

	query := func(idx int) int {
		if val, ok := values[idx]; ok {
			return val
		}
		fmt.Fprintf(writer, "? %d\n", idx)
		writer.Flush()
		var v int
		fmt.Fscan(reader, &v)
		values[idx] = v
		return v
	}

	l, r := 1, n
	for l < r {
		mid := (l + r) / 2
		v1 := query(mid)
		v2 := query(mid + 1)
		if v1 < v2 {
			r = mid
		} else {
			l = mid + 1
		}
	}

	fmt.Fprintf(writer, "! %d\n", l)
	writer.Flush()
}
