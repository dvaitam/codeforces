package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type item struct {
	a int64
	b int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	items := make([]item, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &items[i].a, &items[i].b)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].b < items[j].b })

	l, r := 0, n-1
	var bought, cost int64
	for l <= r {
		for l <= r && items[l].a == 0 {
			l++
		}
		for l <= r && items[r].a == 0 {
			r--
		}
		if l > r {
			break
		}
		if bought >= items[l].b {
			cost += items[l].a
			bought += items[l].a
			items[l].a = 0
			l++
			continue
		}
		need := items[l].b - bought
		if items[r].a <= need {
			cost += items[r].a * 2
			bought += items[r].a
			items[r].a = 0
			r--
		} else {
			cost += need * 2
			bought += need
			items[r].a -= need
		}
	}
	fmt.Fprintln(writer, cost)
}
