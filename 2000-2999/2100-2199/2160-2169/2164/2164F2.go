package main

import (
    "bufio"
    "fmt"
    "os"
)

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    if a < 0 {
        return -a
    }
    return a
}

type segment struct {
    l, r int
}

type dsu struct {
    parent []int
    rank   []int
}

func newDSU(n int) *dsu {
    parent := make([]int, n)
    rank := make([]int, n)
    for i := 0; i < n; i++ {
        parent[i] = i
    }
    return &dsu{parent, rank}
}

func (d *dsu) find(x int) int {
    if d.parent[x] != x {
        d.parent[x] = d.find(d.parent[x])
    }
    return d.parent[x]
}

func (d *dsu) union(x, y int) bool {
    fx, fy := d.find(x), d.find(y)
    if fx == fy {
        return false
    }
    if d.rank[fx] < d.rank[fy] {
        fx, fy = fy, fx
    }
    d.parent[fy] = fx
    if d.rank[fx] == d.rank[fy] {
        d.rank[fx]++
    }
    return true
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n int
    fmt.Fscan(in, &n)
    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &a[i])
    }

    segments := make(map[int][]segment)
    uniqueGcds := make(map[int]struct{})

    for l := 0; l < n; l++ {
        g := 0
        for r := l; r < n; r++ {
            g = gcd(g, a[r])
            segments[g] = append(segments[g], segment{l, r})
            uniqueGcds[g] = struct{}{}
        }
    }

    gcdList := make([]int, 0, len(uniqueGcds))
    for g := range uniqueGcds {
        gcdList = append(gcdList, g)
    }
    // simple insertion sort as len<=something? use sort package (but not imported). We'll do simple bubble? better to implement quick sort.
    quickSort(gcdList, 0, len(gcdList)-1)

    d := newDSU(n)
    var total int64
    edges := 0

    for _, g := range gcdList {
        segs := segments[g]
        for _, seg := range segs {
            if d.union(seg.l, seg.r) {
                total += int64(g)
                edges++
                if edges == n-1 {
                    fmt.Fprintln(out, total)
                    return
                }
            }
        }
    }

    fmt.Fprintln(out, total)
}

func quickSort(arr []int, l, r int) {
    if l >= r {
        return
    }
    pivot := arr[(l+r)/2]
    i, j := l, r
    for i <= j {
        for arr[i] < pivot {
            i++
        }
        for arr[j] > pivot {
            j--
        }
        if i <= j {
            arr[i], arr[j] = arr[j], arr[i]
            i++
            j--
        }
    }
    if l < j {
        quickSort(arr, l, j)
    }
    if i < r {
        quickSort(arr, i, r)
    }
}
