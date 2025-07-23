package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type friend struct {
	money  int64
	factor int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var d int64
	if _, err := fmt.Fscan(in, &n, &d); err != nil {
		return
	}

	friends := make([]friend, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &friends[i].money, &friends[i].factor)
	}

	sort.Slice(friends, func(i, j int) bool {
		return friends[i].money < friends[j].money
	})

	var best int64
	var cur int64
	l := 0
	for r := 0; r < n; r++ {
		cur += friends[r].factor
		for l <= r && friends[r].money-friends[l].money >= d {
			cur -= friends[l].factor
			l++
		}
		if cur > best {
			best = cur
		}
	}

	fmt.Fprintln(out, best)
}
