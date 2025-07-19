package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n, m int
    fmt.Fscan(reader, &n, &m)
    type node struct { val, no int }
    rec := make([]node, m)
    // estimate maximum nodes: columns + group nodes
    maxNodes := m*2 + n*2 + 5
    join := make([]int, maxNodes)
    v := make([][]int, maxNodes)
    cnt := 0
    ans := make([]int, 0, m)

    for i := 0; i < n; i++ {
        for j := 0; j < m; j++ {
            fmt.Fscan(reader, &rec[j].val)
            rec[j].no = j
        }
        sort.Slice(rec, func(i, j int) bool { return rec[i].val < rec[j].val })
        for j := 0; j < m; j++ {
            if rec[j].val < 0 {
                continue
            }
            if j == 0 || rec[j].val != rec[j-1].val {
                cnt++
            }
            // add edge: column -> group end (m+cnt+1)
            end := m + cnt + 1
            v[rec[j].no] = append(v[rec[j].no], end)
            join[end]++
            // add edge: group start (m+cnt) -> column
            start := m + cnt
            v[start] = append(v[start], rec[j].no)
            join[rec[j].no]++
        }
        // separate rows
        cnt++
    }
    totalNodes := m + cnt + 3
    // topological sort
    q := make([]int, 0, totalNodes)
    for i := 0; i < totalNodes && i < len(join); i++ {
        if join[i] == 0 {
            q = append(q, i)
        }
    }
    head := 0
    for head < len(q) {
        t := q[head]
        head++
        if t < m {
            ans = append(ans, t)
        }
        for _, to := range v[t] {
            join[to]--
            if join[to] == 0 {
                q = append(q, to)
            }
        }
    }
    if len(ans) < m {
        fmt.Fprint(writer, -1)
    } else {
        for i := 0; i < m; i++ {
            fmt.Fprint(writer, ans[i]+1)
            if i+1 < m {
                fmt.Fprint(writer, " ")
            }
        }
    }
}
