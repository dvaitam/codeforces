package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for t := 0; t < T; t++ {
		var q int
		fmt.Fscan(reader, &q)
		q1 := 2*q - 1
		items := make([]struct{ a, s, pos int }, q1)
		for i := 0; i < q1; i++ {
			fmt.Fscan(reader, &items[i].a, &items[i].s)
			items[i].pos = i + 1
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].a < items[j].a
		})
		var sum0, sum1 int64
		for i := 0; i < q1; i++ {
			if i&1 == 0 {
				sum0 += int64(items[i].s)
			} else {
				sum1 += int64(items[i].s)
			}
		}
		fmt.Fprint(writer, "YES\n")
		if sum0 >= sum1 {
			for i := 0; i < q1; i += 2 {
				fmt.Fprint(writer, items[i].pos, " ")
			}
		} else {
			for i := 1; i < q1; i += 2 {
				fmt.Fprint(writer, items[i].pos, " ")
			}
			fmt.Fprint(writer, items[q1-1].pos, " ")
		}
		fmt.Fprint(writer, "\n")
	}
}
