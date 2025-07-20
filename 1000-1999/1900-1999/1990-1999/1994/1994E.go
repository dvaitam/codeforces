package main

import (
    "bufio"
    "fmt"
    "os"
    "unicode"
    "container/heap"
)

type FastScanner struct {
    r *bufio.Reader
    buf []byte
}

func NewFastScanner() *FastScanner {
    return &FastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *FastScanner) NextInt() int {
    // Skip non-digit characters
    var c byte
    var err error
    for {
        c, err = fs.r.ReadByte()
        if err != nil {
            return 0
        }
        if unicode.IsDigit(rune(c)) || c == '-' {
            break
        }
    }
    sign := 1
    if c == '-' {
        sign = -1
        c, _ = fs.r.ReadByte()
    }
    res := 0
    for {
        if !unicode.IsDigit(rune(c)) {
            break
        }
        res = res*10 + int(c-'0')
        c, err = fs.r.ReadByte()
        if err != nil {
            break
        }
    }
    return res * sign
}

type IntMaxHeap []int

func (h IntMaxHeap) Len() int           { return len(h) }
func (h IntMaxHeap) Less(i, j int) bool { return h[i] > h[j] } // max-heap
func (h IntMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntMaxHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}
func (h *IntMaxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func main() {
    fs := NewFastScanner()
    out := bufio.NewWriterSize(os.Stdout, 1<<20)
    defer out.Flush()

    t := fs.NextInt()
    for ; t > 0; t-- {
        k := fs.NextInt()
        // capacity slice
        caps := make([]int, 0, k)
        maxN := 0
        for i := 0; i < k; i++ {
            n := fs.NextInt()
            if n > maxN {
                maxN = n
            }
            caps = append(caps, n)
            for j := 2; j <= n; j++ {
                // read and discard parent values
                _ = fs.NextInt()
            }
        }
        h := IntMaxHeap(caps)
        heap.Init(&h)
        res := 0
        // iterate bits from high to low
        for bit := 20; bit >= 0; bit-- {
            val := 1 << bit
            if val > maxN { // no tree individually can hold this bit
                continue
            }
            if h.Len() == 0 {
                break
            }
            top := h[0]
            if top >= val {
                // use this bit
                topCap := heap.Pop(&h).(int)
                topCap -= val
                if topCap > 0 {
                    heap.Push(&h, topCap)
                }
                res |= val
                // We might be able to use the same bit multiple times but it's redundant.
            }
        }
        fmt.Fprintln(out, res)
    }
}


