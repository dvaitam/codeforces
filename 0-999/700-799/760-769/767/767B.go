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

	var ts, tf, t int64
	if _, err := fmt.Fscan(in, &ts, &tf, &t); err != nil {
		return
	}
	var n int
	fmt.Fscan(in, &n)
	arrivals := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arrivals[i])
	}
	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i] < arrivals[j] })

	// Append large sentinel to simplify processing
	arrivals = append(arrivals, int64(1)<<62)
	cur := make([]int64, n+1)
	cur[0] = ts
	for i := 0; i < n; i++ {
		start := cur[i]
		if arrivals[i] > start {
			start = arrivals[i]
		}
		cur[i+1] = start + t
	}

	// try to find zero waiting time
	for i := 0; i <= n; i++ {
		if cur[i]+t <= tf && cur[i] < arrivals[i] {
			fmt.Fprintln(out, cur[i])
			return
		}
	}

	bestTime := ts
	bestWait := int64(1<<63 - 1)
	// candidate at ts
	j := sort.Search(len(arrivals), func(i int) bool { return arrivals[i] > ts })
	start := cur[j]
	if start < ts {
		start = ts
	}
	if start+t <= tf {
		bestWait = start - ts
		bestTime = ts
	}

	// candidates just before each visitor
	for i := 0; i < n; i++ {
		x := arrivals[i] - 1
		if x < 0 {
			continue
		}
		j := sort.Search(len(arrivals), func(k int) bool { return arrivals[k] > x })
		start = cur[j]
		if start < x {
			start = x
		}
		if start+t <= tf {
			wait := start - x
			if wait < bestWait {
				bestWait = wait
				bestTime = x
			}
		}
	}

	if cur[n]+t <= tf && 0 < bestWait {
		bestTime = cur[n]
	}
	fmt.Fprintln(out, bestTime)
}
