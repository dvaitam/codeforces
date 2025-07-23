package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "time"
)

type Fenwick2D struct {
    n, m int
    tree [][]uint64
}

func NewFenwick2D(n, m int) *Fenwick2D {
    tree := make([][]uint64, n+2)
    for i := range tree {
        tree[i] = make([]uint64, m+2)
    }
    return &Fenwick2D{n: n + 1, m: m + 1, tree: tree}
}

func (f *Fenwick2D) add(x, y int, val uint64) {
    for i := x; i <= f.n; i += i & -i {
        row := f.tree[i]
        for j := y; j <= f.m; j += j & -j {
            row[j] ^= val
        }
    }
}

func (f *Fenwick2D) RangeAdd(x1, y1, x2, y2 int, val uint64) {
    f.add(x1, y1, val)
    f.add(x1, y2+1, val)
    f.add(x2+1, y1, val)
    f.add(x2+1, y2+1, val)
}

func (f *Fenwick2D) Query(x, y int) uint64 {
    var res uint64
    for i := x; i > 0; i -= i & -i {
        row := f.tree[i]
        for j := y; j > 0; j -= j & -j {
            res ^= row[j]
        }
    }
    return res
}

type Rect struct{ r1, c1, r2, c2 int }

func main() {
    in := bufio.NewReader(os.Stdin)
    var n, m, q int
    if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
        return
    }
    fw := NewFenwick2D(n+2, m+2)
    rectID := make(map[Rect]uint64)
    rand.Seed(time.Now().UnixNano())

    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    for i := 0; i < q; i++ {
        var t, r1, c1, r2, c2 int
        fmt.Fscan(in, &t, &r1, &c1, &r2, &c2)
        switch t {
        case 1:
            id := rand.Uint64()
            rectID[Rect{r1, c1, r2, c2}] = id
            fw.RangeAdd(r1, c1, r2, c2, id)
        case 2:
            rect := Rect{r1, c1, r2, c2}
            if id, ok := rectID[rect]; ok {
                fw.RangeAdd(r1, c1, r2, c2, id)
                delete(rectID, rect)
            }
        case 3:
            v1 := fw.Query(r1, c1)
            v2 := fw.Query(r2, c2)
            if v1 == v2 {
                fmt.Fprintln(out, "Yes")
            } else {
                fmt.Fprintln(out, "No")
            }
        }
    }
}
