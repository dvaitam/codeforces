package main

import (
    "bufio"
    "fmt"
    "os"
)

// Node represents a vertex of the tree with its children and c value.
type Node struct {
    children []int
    c        int
}

var nodes []Node

// dfs builds the ordering of nodes in the subtree rooted at v.
func dfs(v int) ([]int, bool) {
    order := make([]int, 0)
    for _, ch := range nodes[v].children {
        sub, ok := dfs(ch)
        if !ok {
            return nil, false
        }
        order = append(order, sub...)
    }
    if nodes[v].c > len(order) {
        return nil, false
    }
    idx := nodes[v].c
    order = append(order, 0)
    copy(order[idx+1:], order[idx:])
    order[idx] = v
    return order, true
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }

    nodes = make([]Node, n+1)
    root := 0
    for i := 1; i <= n; i++ {
        var p, c int
        fmt.Fscan(reader, &p, &c)
        nodes[i].c = c
        if p == 0 {
            root = i
        } else {
            nodes[p].children = append(nodes[p].children, i)
        }
    }

    order, ok := dfs(root)
    if !ok || len(order) != n {
        fmt.Fprintln(writer, "NO")
        return
    }

    ans := make([]int, n+1)
    for i, v := range order {
        ans[v] = i + 1
    }

    fmt.Fprintln(writer, "YES")
    for i := 1; i <= n; i++ {
        if i > 1 {
            fmt.Fprint(writer, " ")
        }
        fmt.Fprint(writer, ans[i])
    }
    fmt.Fprintln(writer)
}
