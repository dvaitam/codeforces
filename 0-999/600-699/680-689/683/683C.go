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

	var n1 int
	if _, err := fmt.Fscan(reader, &n1); err != nil {
		return
	}
	set1 := make(map[int]struct{}, n1)
	for i := 0; i < n1; i++ {
		var x int
		fmt.Fscan(reader, &x)
		set1[x] = struct{}{}
	}

	var n2 int
	fmt.Fscan(reader, &n2)
	set2 := make(map[int]struct{}, n2)
	for i := 0; i < n2; i++ {
		var x int
		fmt.Fscan(reader, &x)
		set2[x] = struct{}{}
	}

	var res []int
	for v := range set1 {
		if _, ok := set2[v]; !ok {
			res = append(res, v)
		}
	}
	for v := range set2 {
		if _, ok := set1[v]; !ok {
			res = append(res, v)
		}
	}

	sort.Ints(res)
	fmt.Fprint(writer, len(res))
	for _, v := range res {
		fmt.Fprint(writer, " ", v)
	}
	fmt.Fprintln(writer)
}
