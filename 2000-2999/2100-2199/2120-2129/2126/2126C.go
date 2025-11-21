package main

import (
    "bufio"
    "container/heap"
    "fmt"
    "os"
)

type item struct {
    x     int
    limit int64
}

type maxHeap []item

func (h maxHeap) Len() int { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].limit > h[j].limit }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
    *h = append(*h, x.(item))
}
func (h *maxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func push(h *maxHeap, x int, limit int64) {
    heap.Push(h, item{x, limit})
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n, k int
        fmt.Fscan(in, &n, &k)
        h := make([]int64, n)
        maxH := int64(0)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &h[i])
            if h[i] > maxH {
                maxH = h[i]
            }
        }

        maxPrefix := make([]int64, n)
        maxSuffix := make([]int64, n)
        maxPrefix[0] = h[0]
        for i := 1; i < n; i++ {
            if h[i] > maxPrefix[i-1] {
                maxPrefix[i] = h[i]
            } else {
                maxPrefix[i] = maxPrefix[i-1]
            }
        }
        maxSuffix[n-1] = h[n-1]
        for i := n - 2; i >= 0; i-- {
            if h[i] > maxSuffix[i+1] {
                maxSuffix[i] = h[i]
            } else {
                maxSuffix[i] = maxSuffix[i+1]
            }
        }

        k--
        if h[k] == maxH {
            fmt.Fprintln(out, "YES")
            continue
        }

        hq := &maxHeap{}
        heap.Init(hq)
        push(hq, k, h[k])
        visited := make([]bool, n)
        visited[k] = true
        reachable := false
        for hq.Len() > 0 {
            cur := heap.Pop(hq).(item)
            node := cur.x
            limit := cur.limit
            if h[node] == maxH {
                reachable = true
                break
            }
            for _, nxt := range []int{node - 1, node + 1} {
                if nxt < 0 || nxt >= n || visited[nxt] {
                    continue
                }
                newLimit := limit - absInt64(h[node]-h[nxt])
                if newLimit >= h[nxt] {
                    visited[nxt] = true
                    push(hq, nxt, newLimit)
                }
            }
        }
        if reachable {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}

func absInt64(x int64) int64 {
    if x < 0 {
        return -x
    }
    return x
}
