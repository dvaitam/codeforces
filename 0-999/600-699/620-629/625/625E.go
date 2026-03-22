package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Event struct {
	T    int64
	Type int
	i    int
	j    int
	v_i  int
	v_j  int
}

type PQ []Event

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	if pq[i].T != pq[j].T {
		return pq[i].T < pq[j].T
	}
	return pq[i].Type < pq[j].Type
}
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) {
	*pq = append(*pq, x.(Event))
}
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type Frog struct {
	id        int64
	p         int64
	a         int64
	C         int64
	laps      int64
	next      int
	prev      int
	active    bool
	version   int
	knockouts int64
	orig_idx  int
}

var frogs []Frog
var n int64
var m int64

func Jumps(id int64, n int64, T int64) int64 {
	return (T - id + n) / n
}

func floorDiv(a, b int64) int64 {
	if a >= 0 {
		return a / b
	}
	return (a - b + 1) / b
}

func ceilDiv(a, b int64) int64 {
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b
}

func catchTime(i, j int, T_curr int64) int64 {
	if frogs[i].a == 0 {
		return -1
	}
	R_min := floorDiv(T_curr-frogs[i].id, n) + 2

	delta_j := int64(0)
	if frogs[i].id < frogs[j].id {
		delta_j = -1
	}

	Delta := frogs[i].a - frogs[j].a
	K := frogs[j].C + frogs[i].laps*m - frogs[i].C + delta_j*frogs[j].a

	var R int64
	if Delta > 0 {
		R = ceilDiv(K, Delta)
		if R < R_min {
			R = R_min
		}
	} else {
		if R_min*Delta >= K {
			R = R_min
		} else {
			return -1
		}
	}

	return frogs[i].id + (R-1)*n
}

func updateState(i int, T int64) {
	J := Jumps(frogs[i].id, n, T)
	X := frogs[i].C + J*frogs[i].a
	frogs[i].a -= frogs[i].knockouts
	if frogs[i].a < 0 {
		frogs[i].a = 0
	}
	frogs[i].C = X - J*frogs[i].a
	frogs[i].knockouts = 0
	frogs[i].version++
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanInt64 := func() int64 {
		scanner.Scan()
		val, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		return val
	}

	if !scanner.Scan() {
		return
	}
	n = 0
	nStr := scanner.Text()
	n, _ = strconv.ParseInt(nStr, 10, 64)
	m = scanInt64()

	frogs = make([]Frog, n)
	for i := int64(0); i < n; i++ {
		p := scanInt64()
		a := scanInt64()
		frogs[i] = Frog{
			id:        i + 1,
			p:         p,
			a:         a,
			active:    true,
			orig_idx:  int(i + 1),
		}
	}

	sort.Slice(frogs, func(i, j int) bool {
		return frogs[i].p < frogs[j].p
	})

	for k := 0; k < int(n); k++ {
		frogs[k].prev = (k - 1 + int(n)) % int(n)
		frogs[k].next = (k + 1) % int(n)
		if k == int(n)-1 {
			frogs[k].laps = 1
		} else {
			frogs[k].laps = 0
		}
		frogs[k].C = frogs[k].p
	}

	pq := &PQ{}
	heap.Init(pq)

	for k := 0; k < int(n); k++ {
		T := catchTime(k, frogs[k].next, 0)
		if T != -1 {
			heap.Push(pq, Event{T, 0, k, frogs[k].next, 0, 0})
		}
	}

	activeCount := int(n)

	for pq.Len() > 0 {
		if activeCount <= 1 {
			break
		}
		ev := heap.Pop(pq).(Event)

		if ev.Type == 0 {
			i := ev.i
			j := ev.j
			if !frogs[i].active || !frogs[j].active || frogs[i].next != j {
				continue
			}
			if frogs[i].version != ev.v_i || frogs[j].version != ev.v_j {
				continue
			}

			frogs[j].active = false
			activeCount--
			if activeCount <= 1 {
				break
			}

			nxt := frogs[j].next
			frogs[i].next = nxt
			frogs[nxt].prev = i
			frogs[i].laps += frogs[j].laps
			frogs[i].knockouts++

			T_new := catchTime(i, nxt, ev.T-1)
			if T_new != -1 {
				heap.Push(pq, Event{T_new, 0, i, nxt, frogs[i].version, frogs[nxt].version})
			}

			heap.Push(pq, Event{ev.T, 1, i, 0, frogs[i].version, 0})
		} else {
			i := ev.i
			if !frogs[i].active || frogs[i].version != ev.v_i {
				continue
			}
			if frogs[i].knockouts > 0 {
				updateState(i, ev.T)

				prv := frogs[i].prev
				T_prv := catchTime(prv, i, ev.T)
				if T_prv != -1 {
					heap.Push(pq, Event{T_prv, 0, prv, i, frogs[prv].version, frogs[i].version})
				}

				nxt := frogs[i].next
				T_nxt := catchTime(i, nxt, ev.T)
				if T_nxt != -1 {
					heap.Push(pq, Event{T_nxt, 0, i, nxt, frogs[i].version, frogs[nxt].version})
				}
			}
		}
	}

	fmt.Println(activeCount)
	var out []string
	for _, f := range frogs {
		if f.active {
			out = append(out, strconv.Itoa(f.orig_idx))
		}
	}
	fmt.Println(strings.Join(out, " "))
}
