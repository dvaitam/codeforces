package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	w   int
	val int
}

func query(arr []Pair, w int) int {
	idx := sort.Search(len(arr), func(i int) bool { return arr[i].w >= w })
	if idx == 0 {
		return 0
	}
	return arr[idx-1].val
}

func update(mp *[]Pair, w, val int) {
	arr := *mp
	idx := sort.Search(len(arr), func(i int) bool { return arr[i].w >= w })
	if idx < len(arr) && arr[idx].w == w {
		if arr[idx].val >= val {
			return
		}
		arr[idx].val = val
	} else {
		if idx > 0 && arr[idx-1].val >= val {
			return
		}
		arr = append(arr, Pair{})
		copy(arr[idx+1:], arr[idx:])
		arr[idx] = Pair{w, val}
	}
	j := idx + 1
	for j < len(arr) && arr[j].val <= val {
		j++
	}
	copy(arr[idx+1:], arr[j:])
	arr = arr[:len(arr)-(j-(idx+1))]
	*mp = arr
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	edges := make([]struct{ a, b, w int }, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].a, &edges[i].b, &edges[i].w)
	}
	mp := make([][]Pair, n+1)
	ans := 0
	for _, e := range edges {
		best := query(mp[e.a], e.w)
		val := best + 1
		if val > ans {
			ans = val
		}
		update(&mp[e.b], e.w, val)
	}
	fmt.Println(ans)
}
