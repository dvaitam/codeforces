package main

import (
    "bufio"
    "container/heap"
    "fmt"
    "os"
)

// IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n int
    fmt.Fscan(in, &n)
    ops := make([]struct{ t, x int }, 0, n*2)
    h := &IntHeap{}
    heap.Init(h)

    for i := 0; i < n; i++ {
        var cmd string
        fmt.Fscan(in, &cmd)
        switch cmd {
        case "insert":
            var x int
            fmt.Fscan(in, &x)
            ops = append(ops, struct{ t, x int }{1, x})
            heap.Push(h, x)
        case "getMin":
            var x int
            fmt.Fscan(in, &x)
            if h.Len() == 0 || (*h)[0] > x {
                ops = append(ops, struct{ t, x int }{1, x})
                heap.Push(h, x)
            } else {
                for h.Len() > 0 && (*h)[0] < x {
                    heap.Pop(h)
                    ops = append(ops, struct{ t, x int }{3, 0})
                }
                if h.Len() == 0 || (*h)[0] > x {
                    ops = append(ops, struct{ t, x int }{1, x})
                    heap.Push(h, x)
                }
            }
            ops = append(ops, struct{ t, x int }{2, x})
        case "removeMin":
            if h.Len() == 0 {
                ops = append(ops, struct{ t, x int }{1, 1})
                heap.Push(h, 1)
            }
            heap.Pop(h)
            ops = append(ops, struct{ t, x int }{3, 0})
        }
    }

    fmt.Fprintln(out, len(ops))
    for _, op := range ops {
        switch op.t {
        case 1:
            fmt.Fprintln(out, "insert", op.x)
        case 2:
            fmt.Fprintln(out, "getMin", op.x)
        case 3:
            fmt.Fprintln(out, "removeMin")
        }
    }
}
