package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func tri(x int64) int64 {
	return x * (x + 1) / 2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var x int64
	if _, err := fmt.Fscan(reader, &n, &x); err != nil {
		return
	}
	d := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &d[i])
	}

	arr := make([]int64, 2*n)
	for i := 0; i < 2*n; i++ {
		arr[i] = d[i%n]
	}

	prefDays := make([]int64, 2*n+1)
	prefHugs := make([]int64, 2*n+1)
	for i := 0; i < 2*n; i++ {
		prefDays[i+1] = prefDays[i] + arr[i]
		prefHugs[i+1] = prefHugs[i] + tri(arr[i])
	}

	var best int64
	for r := 1; r <= 2*n; r++ {
		if prefDays[r] < x {
			continue
		}
		need := prefDays[r] - x
		idx := sort.Search(len(prefDays), func(i int) bool { return prefDays[i] > need }) - 1
		if idx < 0 {
			idx = 0
		}
		leftover := need - prefDays[idx]
		partial := tri(arr[idx]) - tri(leftover)
		val := partial + prefHugs[r] - prefHugs[idx+1]
		if val > best {
			best = val
		}
	}

	fmt.Fprintln(writer, best)
}
