package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type contestant struct {
	s    int
	diff int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, S int
	if _, err := fmt.Fscan(in, &n, &S); err != nil {
		return
	}
	people := make([]contestant, n)
	total := 0
	var base int64
	for i := 0; i < n; i++ {
		var si, ai, bi int
		fmt.Fscan(in, &si, &ai, &bi)
		people[i] = contestant{s: si, diff: ai - bi}
		total += si
		base += int64(si) * int64(bi)
	}

	sort.Slice(people, func(i, j int) bool { return people[i].diff > people[j].diff })

	ps := make([]int64, n+1)
	pg := make([]int64, n+1)
	for i := 0; i < n; i++ {
		ps[i+1] = ps[i] + int64(people[i].s)
		pg[i+1] = pg[i] + int64(people[i].s)*int64(people[i].diff)
	}

	pizzas := (total + S - 1) / S
	ans := base
	for x := 0; x <= pizzas; x++ {
		t1 := int64(x * S)
		if t1 > int64(total) {
			t1 = int64(total)
		}
		idx := sort.Search(n, func(i int) bool { return ps[i+1] >= t1 })
		gain := pg[idx]
		if idx < n && t1 > ps[idx] {
			gain += (t1 - ps[idx]) * int64(people[idx].diff)
		}
		if base+gain > ans {
			ans = base + gain
		}
	}
	fmt.Println(ans)
}
