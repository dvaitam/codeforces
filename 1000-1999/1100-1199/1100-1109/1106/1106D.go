package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

// IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n, m int
    if _, err := fmt.Fscan(reader, &n, &m); err != nil {
        return
    }
    graph := make([][]int, n+1)
    for i := 0; i < m; i++ {
        var x, y int
        fmt.Fscan(reader, &x, &y)
        graph[x] = append(graph[x], y)
        graph[y] = append(graph[y], x)
    }
    visited := make([]bool, n+1)
    visited[1] = true
    h := &IntHeap{}
    heap.Init(h)
    for _, v := range graph[1] {
        if !visited[v] {
            visited[v] = true
            heap.Push(h, v)
        }
    }
    writer.WriteString("1 ")
    for h.Len() > 0 {
        u := heap.Pop(h).(int)
        writer.WriteString(strconv.Itoa(u))
        writer.WriteByte(' ')
        for _, v := range graph[u] {
            if !visited[v] {
                visited[v] = true
                heap.Push(h, v)
            }
        }
    }
    writer.WriteByte('\n')
}
