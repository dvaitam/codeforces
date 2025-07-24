package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type State struct {
	time int64
	r1   int64
	r2   int64
	dmg  int
	idx  int
}

type PQ []*State

func (pq PQ) Len() int           { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].time < pq[j].time }
func (pq PQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].idx = i; pq[j].idx = j }
func (pq *PQ) Push(x any) {
	n := len(*pq)
	item := x.(*State)
	item.idx = n
	*pq = append(*pq, item)
}
func (pq *PQ) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var p1, t1, p2, t2, h, s int64
	if _, err := fmt.Fscan(in, &p1, &t1); err != nil {
		return
	}
	fmt.Fscan(in, &p2, &t2)
	fmt.Fscan(in, &h, &s)

	pq := &PQ{}
	heap.Init(pq)
	dist := make(map[[3]int64]int64)

	push := func(time, r1, r2 int64, dmg int) {
		key := [3]int64{r1, r2, int64(dmg)}
		if d, ok := dist[key]; !ok || time < d {
			dist[key] = time
			heap.Push(pq, &State{time: time, r1: r1, r2: r2, dmg: dmg})
		}
	}

	push(0, t1, t2, 0)

	for pq.Len() > 0 {
		st := heap.Pop(pq).(*State)
		key := [3]int64{st.r1, st.r2, int64(st.dmg)}
		if dist[key] != st.time {
			continue
		}
		if int64(st.dmg) >= h {
			fmt.Println(st.time)
			return
		}
		delta := st.r1
		if st.r2 < delta {
			delta = st.r2
		}
		time := st.time + delta
		r1 := st.r1 - delta
		r2 := st.r2 - delta
		if r1 == 0 && r2 == 0 {
			push(time, t1, t2, st.dmg+int(p1+p2-s))
			push(time, t1, 0, st.dmg+int(p1-s))
			push(time, 0, t2, st.dmg+int(p2-s))
			push(time, 0, 0, st.dmg)
		} else if r1 == 0 {
			push(time, t1, r2, st.dmg+int(p1-s))
			push(time+r2, 0, 0, st.dmg)
		} else if r2 == 0 {
			push(time, r1, t2, st.dmg+int(p2-s))
			push(time+r1, 0, 0, st.dmg)
		}
	}
}
