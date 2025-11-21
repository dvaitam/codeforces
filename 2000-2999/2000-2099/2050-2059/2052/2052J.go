package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type task struct {
	d int64
	a int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, q int
		fmt.Fscan(in, &n, &m, &q)

		tasks := make([]task, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &tasks[i].a)
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &tasks[i].d)
		}
		sort.Slice(tasks, func(i, j int) bool { return tasks[i].d < tasks[j].d })
		deadlines := make([]int64, n)
		pref := make([]int64, n+1)
		for i := 0; i < n; i++ {
			deadlines[i] = tasks[i].d
			pref[i+1] = pref[i] + tasks[i].a
		}
		const inf int64 = 1 << 60
		suffMin := make([]int64, n+1)
		suffMin[n] = inf
		for i := n - 1; i >= 0; i-- {
			slack := deadlines[i] - pref[i+1]
			if slack < suffMin[i+1] {
				suffMin[i] = slack
			} else {
				suffMin[i] = suffMin[i+1]
			}
		}

		prefEpisodes := make([]int64, m+1)
		for i := 0; i < m; i++ {
			var l int64
			fmt.Fscan(in, &l)
			prefEpisodes[i+1] = prefEpisodes[i] + l
		}

		answers := make([]int, q)
		for i := 0; i < q; i++ {
			var t int64
			fmt.Fscan(in, &t)
			idx := sort.Search(len(deadlines), func(pos int) bool { return deadlines[pos] > t })
			slack1 := t - pref[idx]
			slack2 := suffMin[idx]
			slack := slack1
			if slack2 < slack {
				slack = slack2
			}
			if slack < 0 {
				slack = 0
			}
			pos := sort.Search(len(prefEpisodes), func(p int) bool { return prefEpisodes[p] > slack }) - 1
			if pos < 0 {
				pos = 0
			}
			answers[i] = pos
		}

		for i := 0; i < q; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, answers[i])
		}
		fmt.Fprintln(out)
	}
}
