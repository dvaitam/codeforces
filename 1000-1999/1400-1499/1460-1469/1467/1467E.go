package main

import (
    "bufio"
    "fmt"
    "os"
)

type state struct {
    node   int
    parent int
    enter  bool
}

func checkRoot(adj [][]int, val []int, root int) bool {
    freq := make(map[int]int)
    stack := []state{{root, -1, true}}
    for len(stack) > 0 {
        s := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if s.enter {
            if freq[val[s.node]] >= 1 {
                return false
            }
            freq[val[s.node]]++
            stack = append(stack, state{s.node, s.parent, false})
            for _, to := range adj[s.node] {
                if to == s.parent {
                    continue
                }
                stack = append(stack, state{to, s.node, true})
            }
        } else {
            freq[val[s.node]]--
        }
    }
    return true
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }
    val := make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Fscan(in, &val[i])
    }
    adj := make([][]int, n+1)
    for i := 0; i < n-1; i++ {
        var u, v int
        fmt.Fscan(in, &u, &v)
        adj[u] = append(adj[u], v)
        adj[v] = append(adj[v], u)
    }
    count := 0
    for r := 1; r <= n; r++ {
        if checkRoot(adj, val, r) {
            count++
        }
    }
    fmt.Println(count)
}

