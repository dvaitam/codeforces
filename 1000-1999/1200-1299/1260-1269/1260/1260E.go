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
	a := make([]int64, n+1)
	friend := -1
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] == -1 {
			friend = i
		}
	}
	if friend == n {
		fmt.Fprintln(writer, 0)
		return
	}

	// r = number of rounds
	r := 0
	for (1 << r) < n {
		r++
	}

	weak := friend - 1
	strong := n - friend

	costs := make([]int64, strong)
	idx := 0
	for i := friend + 1; i <= n; i++ {
		costs[idx] = a[i]
		idx++
	}
	sort.Slice(costs, func(i, j int) bool { return costs[i] < costs[j] })
	prefix := make([]int64, strong+1)
	for i := 0; i < strong; i++ {
		prefix[i+1] = prefix[i] + costs[i]
	}

	best := int64(-1)
	for mask := 0; mask < (1 << r); mask++ {
		// count bribes and simulate
		s := strong
		w := weak
		br := 0
		valid := true
		for i := 0; i < r; i++ {
			if s == 0 {
				if w == 0 {
					valid = false
					break
				}
				w--
				continue
			}
			if mask&(1<<i) != 0 {
				br++
				s = (s - 1 + 1) / 2
			} else {
				if w == 0 {
					valid = false
					break
				}
				w--
				s = (s + 1) / 2
			}
		}
		if !valid || s != 0 || br > strong {
			continue
		}
		if w < 0 {
			continue
		}
		cost := prefix[br]
		if best == -1 || cost < best {
			best = cost
		}
	}

	fmt.Fprintln(writer, best)
}
