package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	h int64
	t int64
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var r int64
	if _, err := fmt.Fscan(in, &n, &r); err != nil {
		return
	}

	arr := make([]interval, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i].h, &arr[i].t)
	}

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].t == arr[j].t {
			return arr[i].h < arr[j].h
		}
		return arr[i].t < arr[j].t
	})

	lastShot := -r // ensures first possible shot is at time 0
	var ans int64

	for _, iv := range arr {
		h := iv.h
		t := iv.t

		if lastShot >= h {
			// already covered by previous shot
			ans++
			continue
		}

		shotTime := max(h, lastShot+r)
		if shotTime < 0 {
			shotTime = 0
		}
		if shotTime <= t {
			lastShot = shotTime
			ans++
		}
	}

	fmt.Fprintln(out, ans)
}
