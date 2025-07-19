package main

import (
    "fmt"
    "os"
)

func main() {
    var n, m int
    if _, err := fmt.Fscan(os.Stdin, &n, &m); err != nil {
        return
    }
    total := n * m
    if total == 1 {
        fmt.Println(1)
        return
    }
    if total <= 3 || (n == 2 && m == 2) {
        fmt.Println(-1)
        return
    }
    a := make([][]int, n)
    for i := range a {
        a[i] = make([]int, m)
    }
    half := total / 2
    for i := 0; i < n; i++ {
        for j := 0; j < m; j++ {
            idx := (i*m + j)/2 + 1
            if (i+j)%2 == 0 {
                a[i][j] = half + idx
            } else {
                a[i][j] = idx
            }
        }
    }
    for i := 0; i < n; i++ {
        for j := 0; j < m; j++ {
            fmt.Fprintf(os.Stdout, "%d ", a[i][j])
        }
        fmt.Fprintln(os.Stdout)
    }
}
