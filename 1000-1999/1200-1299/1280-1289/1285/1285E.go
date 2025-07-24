package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type event struct {
	delta int
	id    int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		l := make([]int64, n)
		r := make([]int64, n)
		coords := make([]int64, 0, 2*n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &l[i], &r[i])
			coords = append(coords, 2*l[i])
			coords = append(coords, 2*r[i]+1)
		}
		sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })
		uniq := make([]int64, 0, len(coords))
		for _, x := range coords {
			if len(uniq) == 0 || uniq[len(uniq)-1] != x {
				uniq = append(uniq, x)
			}
		}
		idx := make(map[int64]int, len(uniq))
		for i, x := range uniq {
			idx[x] = i
		}
		events := make([][]event, len(uniq))
		for i := 0; i < n; i++ {
			events[idx[2*l[i]]] = append(events[idx[2*l[i]]], event{delta: 1, id: i})
			events[idx[2*r[i]+1]] = append(events[idx[2*r[i]+1]], event{delta: -1, id: i})
		}
		active := make(map[int]struct{})
		arr := make([]int, len(uniq)-1)
		single := make([]int, len(uniq)-1)
		for i := range single {
			single[i] = -1
		}
		for j := 0; j < len(uniq)-1; j++ {
			for _, e := range events[j] {
				if e.delta == 1 {
					active[e.id] = struct{}{}
				} else {
					delete(active, e.id)
				}
			}
			arr[j] = len(active)
			if arr[j] == 1 {
				for id := range active {
					single[j] = id
					break
				}
			}
		}
		base := 0
		prev := 0
		for _, c := range arr {
			if prev == 0 && c > 0 {
				base++
			}
			prev = c
		}
		contrib := make([]int, n)
		for j := 0; j < len(arr); j++ {
			if arr[j] == 1 {
				id := single[j]
				leftPos := j > 0 && arr[j-1] > 0
				rightPos := j+1 < len(arr) && arr[j+1] > 0
				if leftPos && rightPos {
					contrib[id]++
				} else if !leftPos && !rightPos {
					contrib[id]--
				}
			}
		}
		ans := 0
		for i := 0; i < n; i++ {
			val := base + contrib[i]
			if val > ans {
				ans = val
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
