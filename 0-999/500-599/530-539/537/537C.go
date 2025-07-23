package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	notes := make([]struct{ d, h int }, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &notes[i].d, &notes[i].h)
	}
	ans := notes[0].h + (notes[0].d - 1)
	for i := 0; i < m-1; i++ {
		d1, h1 := notes[i].d, notes[i].h
		d2, h2 := notes[i+1].d, notes[i+1].h
		dist := d2 - d1
		diff := abs(h2 - h1)
		if diff > dist {
			fmt.Println("IMPOSSIBLE")
			return
		}
		cand := max(h1, h2) + (dist-diff)/2
		if cand > ans {
			ans = cand
		}
	}
	last := notes[m-1]
	if last.h+(n-last.d) > ans {
		ans = last.h + (n - last.d)
	}
	fmt.Println(ans)
}
