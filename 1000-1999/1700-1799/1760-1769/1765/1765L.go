package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Task struct {
	ready int
	pid   int
}

type AvailPQ []Task

func (pq AvailPQ) Len() int            { return len(pq) }
func (pq AvailPQ) Less(i, j int) bool  { return pq[i].pid < pq[j].pid }
func (pq AvailPQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *AvailPQ) Push(x interface{}) { *pq = append(*pq, x.(Task)) }
func (pq *AvailPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

type FuturePQ []Task

func (pq FuturePQ) Len() int { return len(pq) }
func (pq FuturePQ) Less(i, j int) bool {
	if pq[i].ready == pq[j].ready {
		return pq[i].pid < pq[j].pid
	}
	return pq[i].ready < pq[j].ready
}
func (pq FuturePQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *FuturePQ) Push(x interface{}) { *pq = append(*pq, x.(Task)) }
func (pq *FuturePQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

type Event struct {
	day int
	emp int
}

type EventPQ []Event

func (pq EventPQ) Len() int { return len(pq) }
func (pq EventPQ) Less(i, j int) bool {
	if pq[i].day == pq[j].day {
		return pq[i].emp < pq[j].emp
	}
	return pq[i].day < pq[j].day
}
func (pq EventPQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *EventPQ) Push(x interface{}) { *pq = append(*pq, x.(Event)) }
func (pq *EventPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

var holidays []int

func isHoliday(day int) bool {
	idx := sort.SearchInts(holidays, day)
	return idx < len(holidays) && holidays[idx] == day
}

var works [][]bool

func nextWorkDay(emp, day int) int {
	for {
		for !works[emp][(day-1)%7] {
			day++
		}
		if !isHoliday(day) {
			return day
		}
		day++
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	works = make([][]bool, n)
	for i := 0; i < n; i++ {
		var t int
		fmt.Fscan(in, &t)
		arr := make([]bool, 7)
		for j := 0; j < t; j++ {
			var s string
			fmt.Fscan(in, &s)
			var d int
			switch s {
			case "Monday":
				d = 0
			case "Tuesday":
				d = 1
			case "Wednesday":
				d = 2
			case "Thursday":
				d = 3
			case "Friday":
				d = 4
			case "Saturday":
				d = 5
			case "Sunday":
				d = 6
			}
			arr[d] = true
		}
		works[i] = arr
	}

	holidays = make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &holidays[i])
	}

	projects := make([][]int, k)
	positions := make([]int, k)
	totalParts := 0
	for i := 0; i < k; i++ {
		var p int
		fmt.Fscan(in, &p)
		parts := make([]int, p)
		for j := 0; j < p; j++ {
			fmt.Fscan(in, &parts[j])
			parts[j]--
		}
		projects[i] = parts
		totalParts += p
	}

	avail := make([]AvailPQ, n)
	future := make([]FuturePQ, n)
	freeDay := make([]int, n)
	for i := 0; i < n; i++ {
		freeDay[i] = 1
		heap.Init(&avail[i])
		heap.Init(&future[i])
	}

	// initial tasks
	for j := 0; j < k; j++ {
		emp := projects[j][0]
		heap.Push(&avail[emp], Task{ready: 1, pid: j})
	}

	eventPQ := &EventPQ{}
	heap.Init(eventPQ)
	computeEvent := func(e int) int {
		if avail[e].Len() > 0 {
			return nextWorkDay(e, freeDay[e])
		}
		if future[e].Len() > 0 {
			start := freeDay[e]
			if future[e][0].ready > start {
				start = future[e][0].ready
			}
			return nextWorkDay(e, start)
		}
		return 1 << 60
	}

	for e := 0; e < n; e++ {
		if avail[e].Len() > 0 || future[e].Len() > 0 {
			d := computeEvent(e)
			heap.Push(eventPQ, Event{day: d, emp: e})
		}
	}

	results := make([]int, k)
	processed := 0

	for processed < totalParts {
		ev := heap.Pop(eventPQ).(Event)
		e := ev.emp
		expected := computeEvent(e)
		if ev.day != expected {
			heap.Push(eventPQ, Event{day: expected, emp: e})
			continue
		}
		day := ev.day
		for future[e].Len() > 0 && future[e][0].ready <= day {
			t := heap.Pop(&future[e]).(Task)
			heap.Push(&avail[e], t)
		}
		if avail[e].Len() == 0 {
			// no task actually ready, recompute
			nd := computeEvent(e)
			heap.Push(eventPQ, Event{day: nd, emp: e})
			continue
		}
		t := heap.Pop(&avail[e]).(Task)
		pid := t.pid
		positions[pid]++
		processed++
		if positions[pid] == len(projects[pid]) {
			results[pid] = day
		} else {
			nextEmp := projects[pid][positions[pid]]
			ready := day + 1
			if ready <= freeDay[nextEmp] {
				heap.Push(&avail[nextEmp], Task{ready: ready, pid: pid})
			} else {
				heap.Push(&future[nextEmp], Task{ready: ready, pid: pid})
			}
			nd := computeEvent(nextEmp)
			heap.Push(eventPQ, Event{day: nd, emp: nextEmp})
		}
		freeDay[e] = day + 1
		for future[e].Len() > 0 && future[e][0].ready <= freeDay[e] {
			t := heap.Pop(&future[e]).(Task)
			heap.Push(&avail[e], t)
		}
		if avail[e].Len() > 0 || future[e].Len() > 0 {
			nd := computeEvent(e)
			heap.Push(eventPQ, Event{day: nd, emp: e})
		}
	}

	for i := 0; i < k; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, results[i])
	}
	fmt.Fprintln(out)
}
