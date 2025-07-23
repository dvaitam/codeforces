package main

import (
    "bufio"
    "container/heap"
    "fmt"
    "os"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
    old := *h
    x := old[len(old)-1]
    *h = old[:len(old)-1]
    return x
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    fmt.Fscan(reader, &n)
    arr := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &arr[i])
    }
    var q int
    fmt.Fscan(reader, &q)

    for ; q > 0; q-- {
        var l, r int
        fmt.Fscan(reader, &l, &r)
        l--
        r--
        freq := make(map[int]int)
        for i := l; i <= r; i++ {
            freq[arr[i]]++
        }
        if len(freq) == 1 {
            fmt.Fprintln(writer, 0)
            continue
        }
        h := &IntHeap{}
        heap.Init(h)
        for _, f := range freq {
            heap.Push(h, f)
        }
        cost := 0
        for h.Len() > 1 {
            a := heap.Pop(h).(int)
            b := heap.Pop(h).(int)
            sum := a + b
            cost += sum
            heap.Push(h, sum)
        }
        fmt.Fprintln(writer, cost)
    }
}
