package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var tc int
	fmt.Fscan(in, &tc)
	for tc > 0 {
		tc--
		solve(in, out)
	}
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n, s1, s2 int64
	fmt.Fscan(in, &n, &s1, &s2)
	items := make([]struct{ v, id int64 }, n)
	for i := int64(0); i < n; i++ {
		fmt.Fscan(in, &items[i].v)
		items[i].id = i + 1
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].v > items[j].v
	})

	res1 := make([]int64, 0, n)
	res2 := make([]int64, 0, n)
	for _, it := range items {
		// compare cost without overflow: s1*(len(res1)+1) vs s2*(len(res2)+1)
		c1 := s1 * int64(len(res1)+1)
		c2 := s2 * int64(len(res2)+1)
		if c1 < c2 {
			res1 = append(res1, it.id)
		} else if c1 > c2 {
			res2 = append(res2, it.id)
		} else {
			if s1 > s2 {
				res2 = append(res2, it.id)
			} else {
				res1 = append(res1, it.id)
			}
		}
	}
	// output
	fmt.Fprint(out, len(res1))
	for _, id := range res1 {
		fmt.Fprint(out, " ", id)
	}
	fmt.Fprint(out, '\n')
	fmt.Fprint(out, len(res2))
	for _, id := range res2 {
		fmt.Fprint(out, " ", id)
	}
	fmt.Fprint(out, "\n")
}
