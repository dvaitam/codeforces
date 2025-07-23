package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int
	idx int
}

var parent []int
var sz []int
var active []bool
var cnt map[int]int
var segments int

func find(x int) int {
	if parent[x] != x {
		parent[x] = find(parent[x])
	}
	return parent[x]
}

func union(a, b int) {
	ra := find(a)
	rb := find(b)
	if ra == rb {
		return
	}
	la := sz[ra]
	lb := sz[rb]
	cnt[la]--
	if cnt[la] == 0 {
		delete(cnt, la)
	}
	cnt[lb]--
	if cnt[lb] == 0 {
		delete(cnt, lb)
	}
	segments--
	if la < lb {
		ra, rb = rb, ra
		la, lb = lb, la
	}
	parent[rb] = ra
	sz[ra] = la + lb
	cnt[la+lb]++
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	ps := make([]pair, n)
	for i := 0; i < n; i++ {
		ps[i] = pair{arr[i], i}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].val < ps[j].val })

	parent = make([]int, n)
	sz = make([]int, n)
	active = make([]bool, n)
	cnt = make(map[int]int)
	segments = 0

	bestSeg := 0
	bestK := 0

	for _, p := range ps {
		i := p.idx
		active[i] = true
		parent[i] = i
		sz[i] = 1
		cnt[1]++
		segments++
		if i > 0 && active[i-1] {
			union(i, i-1)
		}
		if i+1 < n && active[i+1] {
			union(i, i+1)
		}
		if len(cnt) == 1 {
			var length int
			for k := range cnt {
				length = k
			}
			if cnt[length] == segments {
				k := p.val + 1
				if segments > bestSeg || (segments == bestSeg && k < bestK) {
					bestSeg = segments
					bestK = k
				}
			}
		}
	}

	fmt.Println(bestK)
}
