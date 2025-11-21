package main

import (
    "bufio"
    "container/heap"
    "fmt"
    "math"
    "os"
)

type Edge struct {
    to   int
    bus  int64
    walk int64
}

type Item struct {
    time int64
    node int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].time < pq[j].time }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    *pq = old[:n-1]
    return item
}

var (
    n   int
    t0  int64
    t1  int64
    t2  int64
    adj [][]Edge
)

func can(start int64) bool {
    if start > t0 {
        return false
    }
    const inf int64 = math.MaxInt64 / 4
    dist := make([]int64, n)
    for i := 0; i < n; i++ {
        dist[i] = inf
    }
    dist[0] = start
    pq := &PriorityQueue{}
    heap.Init(pq)
    heap.Push(pq, Item{start, 0})
    for pq.Len() > 0 {
        it := heap.Pop(pq).(Item)
        u := it.node
        cur := it.time
        if cur != dist[u] || cur > t0 {
            continue
        }
        if u == n-1 {
            return true
        }
        for _, e := range adj[u] {
            // walk option
            nt := cur + e.walk
            if nt < dist[e.to] && nt <= t0 {
                dist[e.to] = nt
                heap.Push(pq, Item{nt, e.to})
            }
            // bus option with waiting if needed
            startBus := cur
            if startBus < t1 {
                if startBus+e.bus > t1 {
                    startBus = t2
                }
            } else if startBus < t2 {
                startBus = t2
            }
            if startBus < cur {
                startBus = cur
            }
            nt = startBus + e.bus
            if nt < dist[e.to] && nt <= t0 {
                dist[e.to] = nt
                heap.Push(pq, Item{nt, e.to})
            }
        }
    }
    return false
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var m int
        fmt.Fscan(in, &n, &m)
        fmt.Fscan(in, &t0, &t1, &t2)
        adj = make([][]Edge, n)
        for i := 0; i < m; i++ {
            var u, v int
            var l1, l2 int64
            fmt.Fscan(in, &u, &v, &l1, &l2)
            u--
            v--
            adj[u] = append(adj[u], Edge{v, l1, l2})
            adj[v] = append(adj[v], Edge{u, l1, l2})
        }
        lo, hi := int64(0), t0
        ans := int64(-1)
        for lo <= hi {
            mid := (lo + hi) / 2
            if can(mid) {
                ans = mid
                lo = mid + 1
            } else {
                hi = mid - 1
            }
        }
        fmt.Fprintln(out, ans)
    }
}
