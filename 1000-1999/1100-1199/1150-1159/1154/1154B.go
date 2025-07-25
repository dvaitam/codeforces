package main

import (
    "fmt"
    "sort"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Scan(&a[i])
    }
    sort.Ints(a)
    // collect distinct values
    u := make([]int, 0, n)
    for i, v := range a {
        if i == 0 || v != a[i-1] {
            u = append(u, v)
        }
    }
    m := len(u)
    switch m {
    case 1:
        fmt.Println(0)
    case 2:
        diff := u[1] - u[0]
        if diff%2 == 0 {
            fmt.Println(diff / 2)
        } else {
            fmt.Println(diff)
        }
    case 3:
        d1 := u[1] - u[0]
        d2 := u[2] - u[1]
        if d1 == d2 {
            fmt.Println(d1)
        } else {
            fmt.Println(-1)
        }
    default:
        fmt.Println(-1)
    }
}
