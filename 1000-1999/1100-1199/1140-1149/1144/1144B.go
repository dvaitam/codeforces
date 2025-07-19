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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	ev := make([]int, 0, n)
	od := make([]int, 0, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x%2 == 0 {
			ev = append(ev, x)
		} else {
			od = append(od, x)
		}
	}
	sort.Ints(ev)
	sort.Ints(od)
	diff := len(ev) - len(od)
	sum := 0
	if diff >= 2 {
		// delete smallest diff-1 evens
		for i := 0; i < diff-1; i++ {
			sum += ev[i]
		}
	} else if diff <= -2 {
		// delete smallest (-diff)-1 odds
		k := -diff
		for i := 0; i < k-1; i++ {
			sum += od[i]
		}
	} else {
		fmt.Fprintln(writer, 0)
		return
	}
	fmt.Fprintln(writer, sum)
}
