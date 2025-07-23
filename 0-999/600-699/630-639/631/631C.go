package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Op struct {
	t int
	r int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	t := make([]int, m)
	r := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &t[i], &r[i])
	}

	// Collect relevant operations from the end keeping decreasing r
	ops := make([]Op, 0)
	maxR := 0
	for i := m - 1; i >= 0; i-- {
		if r[i] > maxR {
			ops = append(ops, Op{t: t[i], r: r[i]})
			maxR = r[i]
		}
		if maxR == n {
			break
		}
	}

	if len(ops) == 0 {
		// No manager changes anything
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, arr[i])
		}
		fmt.Fprintln(writer)
		return
	}

	// Reverse operations to process from the first effective manager
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}

	prefix := make([]int, ops[0].r)
	copy(prefix, arr[:ops[0].r])
	sort.Ints(prefix)

	left, right := 0, len(prefix)-1
	for i := 0; i < len(ops); i++ {
		curR := ops[i].r
		nextR := 0
		if i+1 < len(ops) {
			nextR = ops[i+1].r
		}
		if ops[i].t == 1 {
			for pos := curR - 1; pos >= nextR; pos-- {
				arr[pos] = prefix[right]
				right--
			}
		} else {
			for pos := curR - 1; pos >= nextR; pos-- {
				arr[pos] = prefix[left]
				left++
			}
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, arr[i])
	}
	fmt.Fprintln(writer)
}
