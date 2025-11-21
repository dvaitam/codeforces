package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)
        pos := make(map[int][]int)
        b := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &b[i])
            pos[b[i]] = append(pos[b[i]], i)
        }
        keys := make([]int, 0, len(pos))
        for k := range pos {
            keys = append(keys, k)
        }
        sort.Ints(keys)
        res := make([]int, n)
        label := 1
        possible := true
        for _, v := range keys {
            idxs := pos[v]
            if len(idxs)%v != 0 {
                possible = false
                break
            }
            for i := 0; i < len(idxs); i += v {
                for j := 0; j < v; j++ {
                    res[idxs[i+j]] = label
                }
                label++
            }
        }
        if !possible {
            fmt.Fprintln(out, -1)
            continue
        }
        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, res[i])
        }
        fmt.Fprintln(out)
    }
}
