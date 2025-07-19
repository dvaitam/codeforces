package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func nextInt() int {
    var x int
    fmt.Fscan(reader, &x)
    return x
}

type node struct {
    fi int
    se int
}

func solve() (bool, [][2]int) {
    n := nextInt()
    c := make([]int, n+2)
    d := make([]node, n+2)
    var su int
    for i := 1; i <= n; i++ {
        c[i] = nextInt()
        d[i] = node{fi: c[i], se: i}
        su += c[i]
    }
    if su != 2*n {
        return false, nil
    }
    // sort by fi then se on d[1..n]
    ds := d[1 : n+1]
    sort.Slice(ds, func(i, j int) bool {
        if ds[i].fi != ds[j].fi {
            return ds[i].fi < ds[j].fi
        }
        return ds[i].se < ds[j].se
    })
    ans := make([][2]int, 0, 2*n)
    add := func(x, y int) {
        ans = append(ans, [2]int{d[x].se, d[y].se})
        d[x].fi--
        d[y].fi--
    }
    // case cycle
    if d[1].fi == 2 {
        for i := 1; i <= n; i++ {
            u := d[i].se
            v := d[i%n+1].se
            ans = append(ans, [2]int{u, v})
        }
        return true, ans
    }
    if d[2].fi != 1 || d[n-1].fi <= 2 {
        return false, nil
    }
    top, s := 1, 1
    for top <= n && d[top].fi == 1 {
        top++
    }
    // single big block >2
    if d[top].fi > 2 {
        add(top, n)
        for i := top; i < n; i++ {
            add(i, i+1)
        }
        for i := top; i <= n; i++ {
            for d[i].fi > 0 {
                add(s, i)
                s++
            }
        }
        return true, ans
    }
    // complex cases
    for {
        if top+1 == n {
            add(top-2, top)
            add(top-1, top+1)
            add(top, top+1)
            add(top, top+1)
            for d[top].fi > 0 {
                add(s, top)
                s++
            }
            for d[top+1].fi > 0 {
                add(s, top+1)
                s++
            }
            break
        }
        if top+2 == n {
            if c[d[top-1].se] == 1 {
                return false, nil
            }
            if d[top+2].fi <= 3 {
                return false, nil
            }
            add(top, top+2)
            for d[top].fi > 0 {
                add(s, top)
                s++
            }
            // swap d[top] and d[top-2]
            d[top], d[top-2] = d[top-2], d[top]
            top++
            continue
        }
        if c[d[top-1].se] != 1 && top+4 == n && d[top+2].fi >= 3 {
            add(top, top+2)
            for d[top].fi > 0 {
                add(s, top)
                s++
            }
            add(top-2, top+1)
            add(top-1, top+2)
            for d[top+1].fi > 1 {
                add(s, top+1)
                s++
            }
            for d[top+2].fi > 1 {
                add(s, top+2)
                s++
            }
            top += 3
            continue
        }
        add(top-2, top)
        add(top-1, top+1)
        for d[top].fi > 1 {
            add(s, top)
            s++
        }
        for d[top+1].fi > 1 {
            add(s, top+1)
            s++
        }
        top += 2
    }
    return true, ans
}

func main() {
    defer writer.Flush()
    TT := nextInt()
    for t := 0; t < TT; t++ {
        ok, ans := solve()
        if ok {
            fmt.Fprintln(writer, "Yes")
            for _, e := range ans {
                fmt.Fprintln(writer, e[0], e[1])
            }
        } else {
            fmt.Fprintln(writer, "No")
        }
    }
}
