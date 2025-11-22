package main

import (
    "bufio"
    "fmt"
    "os"
)

type Fenwick struct {
    n   int
    bit []int
}

func NewFenwick(n int) *Fenwick {
    return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(idx, delta int) {
    for idx <= f.n {
        f.bit[idx] += delta
        idx += idx & -idx
    }
}

func (f *Fenwick) Sum(idx int) int {
    res := 0
    for idx > 0 {
        res += f.bit[idx]
        idx -= idx & -idx
    }
    return res
}

// LowerBound returns the smallest index such that prefix sum >= k (k >= 1).
func (f *Fenwick) LowerBound(k int) int {
    idx := 0
    bitMask := 1
    for bitMask<<1 <= f.n {
        bitMask <<= 1
    }
    for bitMask > 0 {
        next := idx + bitMask
        if next <= f.n && f.bit[next] < k {
            k -= f.bit[next]
            idx = next
        }
        bitMask >>= 1
    }
    return idx + 1
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var q int
    fmt.Fscan(in, &q)
    const MAX = 300000
    fw := NewFenwick(MAX)
    present := make([]bool, MAX+1)

    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    for ; q > 0; q-- {
        var tp int
        fmt.Fscan(in, &tp)
        if tp == 1 {
            var s int
            fmt.Fscan(in, &s)
            if present[s] {
                present[s] = false
                fw.Add(s, -1)
            } else {
                present[s] = true
                fw.Add(s, 1)
            }
        } else if tp == 2 {
            var k int
            fmt.Fscan(in, &k)
            xorVal := 0
            for i := 0; i < k; i++ {
                var a int
                fmt.Fscan(in, &a)
                cnt := fw.Sum(a)
                prev := 0
                if cnt > 0 {
                    prev = fw.LowerBound(cnt)
                }
                xorVal ^= (a - prev)
            }
            if xorVal != 0 {
                fmt.Fprintln(out, "First")
            } else {
                fmt.Fprintln(out, "Second")
            }
        }
    }
}

