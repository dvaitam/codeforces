package main

import (
    "bufio"
    "fmt"
    "os"
)

func solve(s string) int {
    bpos := []int{}
    for i := 0; i < len(s); i++ {
        if s[i] == 'B' {
            bpos = append(bpos, i)
        }
    }
    if len(bpos) == 0 {
        return 0
    }
    k := len(bpos)
    segments := make([]int, k+1)
    segments[0] = bpos[0]
    for i := 0; i < k-1; i++ {
        segments[i+1] = bpos[i+1] - bpos[i] - 1
    }
    segments[k] = len(s) - 1 - bpos[k-1]
    edges := make([]int, 0, 2*k)
    for i := 0; i < k; i++ {
        edges = append(edges, segments[i])
        edges = append(edges, segments[i+1])
    }
    if len(edges) == 0 {
        return 0
    }
    dpPrev2 := 0
    dpPrev1 := max(0, edges[0])
    for i := 2; i <= len(edges); i++ {
        val := max(dpPrev1, dpPrev2+edges[i-1])
        dpPrev2, dpPrev1 = dpPrev1, val
    }
    return dpPrev1
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var s string
        fmt.Fscan(reader, &s)
        fmt.Fprintln(writer, solve(s))
    }
}
