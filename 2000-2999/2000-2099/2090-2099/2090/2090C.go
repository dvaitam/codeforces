package main

import (
    "bufio"
    "container/heap"
    "fmt"
    "os"
)

type cell struct {
    x, y int
    dist int
}

type state struct {
    typ int
    s   int
    a   int
    x   int
    y   int
    dist int
}

type priority []*state

func (p priority) Len() int { return len(p) }
func (p priority) Less(i, j int) bool {
    if p[i].dist != p[j].dist {
        return p[i].dist < p[j].dist
    }
    if p[i].x != p[j].x {
        return p[i].x < p[j].x
    }
    return p[i].y < p[j].y
}
func (p priority) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p *priority) Push(v interface{}) {
    *p = append(*p, v.(*state))
}
func (p *priority) Pop() interface{} {
    old := *p
    v := old[len(old)-1]
    *p = old[:len(old)-1]
    return v
}

func makeState(typ, s, a int) *state {
    var x, y, dist int
    switch typ {
    case 0: // (1,1)
        x = 3*a + 1
        y = 3*(s-a) + 1
        dist = 3*s + 2
    case 1: // (1,2)
        x = 3*a + 1
        y = 3*(s-a) + 2
        dist = 3*s + 3
    case 2: // (2,1)
        x = 3*a + 2
        y = 3*(s-a) + 1
        dist = 3*s + 3
    case 3: // (2,2)
        x = 3*a + 2
        y = 3*(s-a) + 2
        dist = 3*s + 6
    }
    return &state{typ: typ, s: s, a: a, x: x, y: y, dist: dist}
}

func generateCells(limit int) ([]cell, []cell) {
    // returns: full cell order, and type-0 cells (one per table) order
    pq := priority{}
    for t := 0; t < 4; t++ {
        heap.Push(&pq, makeState(t, 0, 0))
    }
    all := make([]cell, 0, limit)
    tables := make([]cell, 0, limit)
    for len(all) < limit {
        cur := heap.Pop(&pq).(*state)
        all = append(all, cell{cur.x, cur.y, cur.dist})
        if cur.typ == 0 {
            tables = append(tables, cell{cur.x, cur.y, cur.dist})
        }
        // advance generator
        if cur.a+1 <= cur.s {
            heap.Push(&pq, makeState(cur.typ, cur.s, cur.a+1))
        } else {
            heap.Push(&pq, makeState(cur.typ, cur.s+1, 0))
        }
    }
    return all, tables
}

func key(x, y int) int64 {
    return (int64(x) << 32) ^ int64(uint32(y))
}

func tableKey(x, y int) int64 {
    // table coordinates (a,b) derived from top-left cell (3a+1,3b+1)
    a := (x - 1) / 3
    b := (y - 1) / 3
    return (int64(a) << 32) ^ int64(uint32(b))
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var q int
    if _, err := fmt.Fscan(in, &q); err != nil {
        return
    }
    type test struct {
        n int
        t []int
    }
    tests := make([]test, q)
    total := 0
    for i := 0; i < q; i++ {
        var n int
        fmt.Fscan(in, &n)
        arr := make([]int, n)
        for j := 0; j < n; j++ {
            fmt.Fscan(in, &arr[j])
        }
        tests[i] = test{n: n, t: arr}
        total += n
    }

    // generate enough cells and tables
    limit := total*2 + 100
    cells, tables := generateCells(limit)

    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    for _, tc := range tests {
        occupied := make(map[int64]bool)
        tableOcc := make(map[int64]bool)
        ci, ti := 0, 0
        for _, typ := range tc.t {
            if typ == 1 {
                // nearest vacant cell
                for occupied[key(cells[ci].x, cells[ci].y)] {
                    ci++
                }
                c := cells[ci]
                fmt.Fprintln(out, c.x, c.y)
                occupied[key(c.x, c.y)] = true
                tkey := tableKey(c.x, c.y)
                tableOcc[tkey] = true
                ci++
            } else {
                // nearest completely empty table -> cell (3a+1,3b+1)
                for tableOcc[tableKey(tables[ti].x, tables[ti].y)] || occupied[key(tables[ti].x, tables[ti].y)] {
                    ti++
                }
                c := tables[ti]
                fmt.Fprintln(out, c.x, c.y)
                occupied[key(c.x, c.y)] = true
                tableOcc[tableKey(c.x, c.y)] = true
                ti++
            }
        }
    }
}
