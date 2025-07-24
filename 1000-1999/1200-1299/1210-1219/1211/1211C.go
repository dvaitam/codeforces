package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type day struct {
	a int64
	b int64
	c int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	days := make([]day, n)
	var minNeed, maxPossible int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &days[i].a, &days[i].b, &days[i].c)
		minNeed += days[i].a
		maxPossible += days[i].b
	}
	if k < minNeed || k > maxPossible {
		fmt.Fprintln(out, -1)
		return
	}

	totalCost := int64(0)
	for i := 0; i < n; i++ {
		totalCost += days[i].a * days[i].c
		days[i].b -= days[i].a
	}
	remaining := k - minNeed
	sort.Slice(days, func(i, j int) bool { return days[i].c < days[j].c })
	for _, d := range days {
		if remaining == 0 {
			break
		}
		take := d.b
		if take > remaining {
			take = remaining
		}
		totalCost += take * d.c
		remaining -= take
	}
	fmt.Fprintln(out, totalCost)
}
