package main

import (
    "bufio"
    "container/heap"
    "fmt"
    "math"
    "os"
    "sort"
)

type busInfo struct {
    s, t int64
    time float64
}

type person struct {
    pos int64
    idx int
}

type heapEntry struct {
    end  int64
    time float64
}

type busHeap []heapEntry

func (h busHeap) Len() int            { return len(h) }
func (h busHeap) Less(i, j int) bool  { return h[i].time < h[j].time }
func (h busHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *busHeap) Push(x interface{}) { *h = append(*h, x.(heapEntry)) }
func (h *busHeap) Pop() interface{} {
    old := *h
    n := len(old)
    item := old[n-1]
    *h = old[:n-1]
    return item
}

func main() {
    in := bufio.NewReader(os.Stdin)

    var n, m int
    var L, x, y int64
    if _, err := fmt.Fscan(in, &n, &m, &L, &x, &y); err != nil {
        return
    }

    buses := make([]busInfo, n)
    for i := 0; i < n; i++ {
        var s, t int64
        fmt.Fscan(in, &s, &t)
        info := busInfo{s: s, t: t}
        if t >= L {
            info.time = float64(L-s) / float64(x)
        } else {
            ride := float64(t-s) / float64(x)
            walk := float64(L-t) / float64(y)
            info.time = ride + walk
        }
        buses[i] = info
    }

    persons := make([]person, m)
    answers := make([]float64, m)
    for i := 0; i < m; i++ {
        var p int64
        fmt.Fscan(in, &p)
        persons[i] = person{pos: p, idx: i}
    }

    sort.Slice(buses, func(i, j int) bool {
        if buses[i].s == buses[j].s {
            return buses[i].t < buses[j].t
        }
        return buses[i].s < buses[j].s
    })

    sort.Slice(persons, func(i, j int) bool {
        return persons[i].pos < persons[j].pos
    })

    h := &busHeap{}
    heap.Init(h)
    busPtr := 0

    for _, p := range persons {
        for busPtr < n && buses[busPtr].s <= p.pos {
            heap.Push(h, heapEntry{end: buses[busPtr].t, time: buses[busPtr].time})
            busPtr++
        }
        for h.Len() > 0 && (*h)[0].end < p.pos {
            heap.Pop(h)
        }

        best := math.Inf(1)
        if h.Len() > 0 {
            best = (*h)[0].time
        }

        walk := float64(L-p.pos) / float64(y)
        if walk < 0 {
            walk = 0
        }
        if best > walk {
            answers[p.idx] = walk
        } else {
            answers[p.idx] = best
        }
    }

    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    for i := 0; i < m; i++ {
        fmt.Fprintf(out, "%.10f\n", answers[i])
    }
}
