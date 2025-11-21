package main

import (
    "bufio"
    "fmt"
    "os"
)

type edge struct {
    u, v int
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)
        vols := make([]int, n)
        pits := make([]int, n)
        volMap := make(map[int]int)
        pitMap := make(map[int]int)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &vols[i], &pits[i])
            if _, ok := volMap[vols[i]]; !ok {
                volMap[vols[i]] = len(volMap)
            }
            if _, ok := pitMap[pits[i]]; !ok {
                pitMap[pits[i]] = len(pitMap)
            }
        }
        volCnt := len(volMap)
        pitCnt := len(pitMap)
        tot := volCnt + pitCnt

        edges := make([]edge, n)
        adj := make([][]int, tot)
        deg := make([]int, tot)
        for i := 0; i < n; i++ {
            u := volMap[vols[i]]
            v := pitMap[pits[i]] + volCnt
            edges[i] = edge{u, v}
            adj[u] = append(adj[u], i)
            adj[v] = append(adj[v], i)
            deg[u]++
            deg[v]++
        }
        if n == 0 {
            fmt.Fprintln(out, "NO")
            continue
        }
        // connectedness check
        visited := make([]bool, tot)
        stack := []int{edges[0].u}
        visited[edges[0].u] = true
        for len(stack) > 0 {
            v := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            for _, eid := range adj[v] {
                to := edges[eid].u
                if to == v {
                    to = edges[eid].v
                }
                if !visited[to] {
                    visited[to] = true
                    stack = append(stack, to)
                }
            }
        }
        connected := true
        for v := 0; v < tot; v++ {
            if deg[v] > 0 && !visited[v] {
                connected = false
                break
            }
        }
        if !connected {
            fmt.Fprintln(out, "NO")
            continue
        }
        // Euler trail check
        odd := []int{}
        for v := 0; v < tot; v++ {
            if deg[v]%2 == 1 {
                odd = append(odd, v)
            }
        }
        if len(odd) != 0 && len(odd) != 2 {
            fmt.Fprintln(out, "NO")
            continue
        }
        start := edges[0].u
        if len(odd) == 2 {
            start = odd[0]
        }
        used := make([]bool, n)
        pos := make([]int, tot)
        stackV := []int{start}
        stackE := []int{}
        order := make([]int, 0, n)

        for len(stackV) > 0 {
            v := stackV[len(stackV)-1]
            for pos[v] < len(adj[v]) && used[adj[v][pos[v]]] {
                pos[v]++
            }
            if pos[v] == len(adj[v]) {
                stackV = stackV[:len(stackV)-1]
                if len(stackE) > 0 {
                    eid := stackE[len(stackE)-1]
                    stackE = stackE[:len(stackE)-1]
                    order = append(order, eid)
                }
            } else {
                eid := adj[v][pos[v]]
                pos[v]++
                if used[eid] {
                    continue
                }
                used[eid] = true
                stackE = append(stackE, eid)
                to := edges[eid].u
                if to == v {
                    to = edges[eid].v
                }
                stackV = append(stackV, to)
            }
        }
        if len(order) != n {
            fmt.Fprintln(out, "NO")
            continue
        }
        // reverse order
        for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
            order[i], order[j] = order[j], order[i]
        }
        fmt.Fprintln(out, "YES")
        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, order[i]+1)
        }
        fmt.Fprintln(out)
    }
}
