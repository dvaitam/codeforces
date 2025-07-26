package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)
	var m int
	fmt.Fscan(reader, &m)

	pairCnt := 26 * 26
	pairID := make([]int, n-1)
	posByPair := make([][]int, pairCnt)
	for i := 0; i < n-1; i++ {
		p := int(s[i]-'a')*26 + int(s[i+1]-'a')
		pairID[i] = p
		posByPair[p] = append(posByPair[p], i)
	}

	INF := int(1e9)

	// nearest distance from each index to nearest occurrence of pair p
	nearest := make([][]int, pairCnt)
	for p := 0; p < pairCnt; p++ {
		arr := make([]int, n-1)
		for i := range arr {
			arr[i] = INF
		}
		for _, idx := range posByPair[p] {
			arr[idx] = 0
		}
		last := -INF
		for i := 0; i < n-1; i++ {
			if arr[i] == 0 {
				last = i
			}
			if last != -INF {
				d := i - last
				if d < arr[i] {
					arr[i] = d
				}
			}
		}
		last = INF
		for i := n - 2; i >= 0; i-- {
			if arr[i] == 0 {
				last = i
			}
			if last != INF {
				d := last - i
				if d < arr[i] {
					arr[i] = d
				}
			}
		}
		nearest[p] = arr
	}

	pairStart := n - 1
	totalNodes := pairStart + pairCnt

	// Deque implementation
	type Deque struct {
		data []int
		l, r int
	}
	pushFront := func(d *Deque, x int) {
		d.l = (d.l - 1 + len(d.data)) % len(d.data)
		d.data[d.l] = x
	}
	pushBack := func(d *Deque, x int) {
		d.data[d.r] = x
		d.r = (d.r + 1) % len(d.data)
	}
	popFront := func(d *Deque) int {
		x := d.data[d.l]
		d.l = (d.l + 1) % len(d.data)
		return x
	}
	empty := func(d *Deque) bool { return d.l == d.r }

	bfs := func(p int) []int {
		dist := make([]int, totalNodes)
		for i := range dist {
			dist[i] = INF
		}
		dq := &Deque{data: make([]int, totalNodes*2+5)}
		start := pairStart + p
		dist[start] = 0
		pushFront(dq, start)
		for !empty(dq) {
			v := popFront(dq)
			dval := dist[v]
			if v >= pairStart {
				idx := v - pairStart
				for _, pos := range posByPair[idx] {
					if dist[pos] > dval {
						dist[pos] = dval
						pushFront(dq, pos)
					}
				}
			} else {
				if v > 0 && dist[v-1] > dval+1 {
					dist[v-1] = dval + 1
					pushBack(dq, v-1)
				}
				if v+1 < n-1 && dist[v+1] > dval+1 {
					dist[v+1] = dval + 1
					pushBack(dq, v+1)
				}
				pn := pairStart + pairID[v]
				if dist[pn] > dval+1 {
					dist[pn] = dval + 1
					pushBack(dq, pn)
				}
			}
		}
		return dist[:n-1]
	}

	distFromPair := make([][]int, pairCnt)
	for p := 0; p < pairCnt; p++ {
		distFromPair[p] = bfs(p)
	}

	for ; m > 0; m-- {
		var f, t int
		fmt.Fscan(reader, &f, &t)
		f--
		t--
		ans := abs(f - t)
		for p := 0; p < pairCnt; p++ {
			nf := nearest[p][f]
			if nf >= INF {
				continue
			}
			dt := distFromPair[p][t]
			if dt >= INF {
				continue
			}
			cost := nf + 1 + dt
			if cost < ans {
				ans = cost
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
