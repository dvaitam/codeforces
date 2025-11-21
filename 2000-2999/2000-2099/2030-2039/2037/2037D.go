package main

import (
    "bufio"
    "container/heap"
    "fmt"
    "os"
)

type MaxHeap []int64

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MaxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n, m int
        var L int64
        fmt.Fscan(in, &n, &m, &L)
        hurdles := make([][2]int64, n)
        for i := 0; i < n; i++ {
            var l, r int64
            fmt.Fscan(in, &l, &r)
            hurdles[i][0] = l
            hurdles[i][1] = r
        }
        powerPos := make([]int64, m)
        powerVal := make([]int64, m)
        for i := 0; i < m; i++ {
            fmt.Fscan(in, &powerPos[i], &powerVal[i])
        }

        if n == 0 {
            fmt.Fprintln(out, 0)
            continue
        }

        curPower := int64(1)
        used := int64(0)
        h := &MaxHeap{}
        heap.Init(h)
        idx := 0
        possible := true
        for _, hurdle := range hurdles {
            l := hurdle[0]
            r := hurdle[1]
            for idx < m && powerPos[idx] < l {
                heap.Push(h, powerVal[idx])
                idx++
            }
            need := (r - l + 2)
            for curPower < need {
                if h.Len() == 0 {
                    possible = false
                    break
                }
                curPower += heap.Pop(h).(int64)
                used++
            }
            if !possible {
                break
            }
        }
        if possible {
            fmt.Fprintln(out, used)
        } else {
            fmt.Fprintln(out, -1)
        }
    }
}
