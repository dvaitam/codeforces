package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type event struct {
	t int
	d int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	events := make([]event, 0, 2*n)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		events = append(events, event{t: l, d: 1})
		events = append(events, event{t: r, d: -1})
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].t == events[j].t {
			return events[i].d > events[j].d
		}
		return events[i].t < events[j].t
	})

	cnt := 0
	for _, e := range events {
		cnt += e.d
		if cnt > 2 {
			fmt.Fprintln(out, "NO")
			return
		}
	}
	fmt.Fprintln(out, "YES")
}
