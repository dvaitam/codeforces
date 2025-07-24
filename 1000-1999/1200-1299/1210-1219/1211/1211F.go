package main

import (
    "bufio"
    "fmt"
    "os"
)

type Edge struct {
    to  int
    idx int
}

func offsetFor(s string) int {
    pat := "kotlin"
    n := len(s)
    for start := 0; start < 6; start++ {
        ok := true
        for i := 0; i < n; i++ {
            if s[i] != pat[(start+i)%6] {
                ok = false
                break
            }
        }
        if ok {
            return start
        }
    }
    return 0
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }
    pieces := make([]string, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &pieces[i])
    }

    adj := make([][]Edge, 6)
    for i, s := range pieces {
        start := offsetFor(s)
        end := (start + len(s)) % 6
        adj[start] = append(adj[start], Edge{to: end, idx: i})
    }

    iter := make([]int, 6)
    stack := []int{0}
    estack := []int{}
    order := []int{}

    for len(stack) > 0 {
        v := stack[len(stack)-1]
        if iter[v] < len(adj[v]) {
            e := adj[v][iter[v]]
            iter[v]++
            stack = append(stack, e.to)
            estack = append(estack, e.idx)
        } else {
            stack = stack[:len(stack)-1]
            if len(estack) > 0 {
                order = append(order, estack[len(estack)-1])
                estack = estack[:len(estack)-1]
            }
        }
    }

    for i := len(order) - 1; i >= 0; i-- {
        if i != len(order)-1 {
            fmt.Fprint(out, " ")
        }
        fmt.Fprint(out, order[i]+1)
    }
    if len(order) > 0 {
        fmt.Fprintln(out)
    }
}
