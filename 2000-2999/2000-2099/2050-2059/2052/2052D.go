package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }

    typeNode := make([]int, n+1)
	stCount, utCount := 0, 0

    for i := 1; i <= n; i++ {
        var op, res string
        fmt.Fscan(in, &op, &res)
        switch op {
        case "set":
            if res == "true" {
				typeNode[i] = 1
				stCount++
            } else {
                typeNode[i] = 2
            }
        case "unset":
            if res == "true" {
				typeNode[i] = 3
				utCount++
            } else {
                typeNode[i] = 0
            }
        }
    }

    if stCount > 1 || utCount > 1 {
        fmt.Fprintln(out, -1)
        return
    }
    if utCount == 1 && stCount == 0 {
        fmt.Fprintln(out, -1)
        return
    }
    if stCount == 0 {
        for i := 1; i <= n; i++ {
            if typeNode[i] != 0 {
                fmt.Fprintln(out, -1)
                return
            }
        }
    }

    var m int
    fmt.Fscan(in, &m)
    adj := make([][]int, n+1)
    indeg := make([]int, n+1)
    for i := 0; i < m; i++ {
        var a, b int
        fmt.Fscan(in, &a, &b)
        adj[a] = append(adj[a], b)
        indeg[b]++
    }

    queues := make([][]int, 4)
    heads := make([]int, 4)
    push := func(t, v int) {
        queues[t] = append(queues[t], v)
    }

    for i := 1; i <= n; i++ {
        if indeg[i] == 0 {
            push(typeNode[i], i)
        }
    }

    pop := func(t int) (int, bool) {
        if heads[t] >= len(queues[t]) {
            return 0, false
        }
        v := queues[t][heads[t]]
        heads[t]++
        return v, true
    }

    order := make([]int, 0, n)
    processed := 0
    stateTrue := false
    usedST := stCount == 0
    usedUT := utCount == 0

    for processed < n {
        var v int
        var ok bool
        if !stateTrue {
            if v, ok = pop(0); ok {
                // nothing
            } else if !usedST {
                if v, ok = pop(1); !ok {
                    fmt.Fprintln(out, -1)
                    return
                }
                usedST = true
                stateTrue = true
            } else {
                fmt.Fprintln(out, -1)
                return
            }
        } else {
            if v, ok = pop(2); ok {
                // nothing
            } else if !usedUT {
                if v, ok = pop(3); !ok {
                    fmt.Fprintln(out, -1)
                    return
                }
                usedUT = true
                stateTrue = false
            } else {
                fmt.Fprintln(out, -1)
                return
            }
        }

        order = append(order, v)
        processed++
        for _, nb := range adj[v] {
            indeg[nb]--
            if indeg[nb] == 0 {
                push(typeNode[nb], nb)
            }
        }
    }

    // Verify final state: if state true but UT nonexistent -> OK? stateTrue indicates last state after last op. no extra check needed.
    for i, v := range order {
        if i > 0 {
            fmt.Fprint(out, " ")
        }
        fmt.Fprint(out, v)
    }
    fmt.Fprintln(out)
}
