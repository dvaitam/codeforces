package main

import (
    "bufio"
    "fmt"
    "os"
)

type point struct{r, c int}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, m, x int
    fmt.Fscan(in, &n, &m, &x)

    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &a[i])
    }
    b := make([]int, m)
    for j := 0; j < m; j++ {
        fmt.Fscan(in, &b[j])
    }

    good := make(map[point]struct{})
    for i := 0; i < n; i++ {
        for j := 0; j < m; j++ {
            if a[i]+b[j] <= x {
                good[point{i, j}] = struct{}{}
            }
        }
    }

    visited := make(map[point]bool)
    dirs := []point{{1,0},{-1,0},{0,1},{0,-1}}
    var stack []point
    comp := 0
    for p := range good {
        if visited[p] {
            continue
        }
        comp++
        stack = append(stack[:0], p)
        visited[p] = true
        for len(stack) > 0 {
            cur := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            for _, d := range dirs {
                np := point{cur.r+d.r, cur.c+d.c}
                if _, ok := good[np]; ok && !visited[np] {
                    visited[np] = true
                    stack = append(stack, np)
                }
            }
        }
    }

    fmt.Fprintln(out, comp)
}

