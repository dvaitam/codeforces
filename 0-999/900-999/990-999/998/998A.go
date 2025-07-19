package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    a := make([]int, n)
    sum := 0
    for i := 0; i < n; i++ {
        fmt.Scan(&a[i])
        sum += a[i]
    }
    if n == 1 || (n == 2 && a[0] == a[1]) {
        fmt.Println(-1)
        return
    }
    if sum-a[0] != a[0] {
        fmt.Println(1)
        fmt.Println(1)
    } else {
        fmt.Println(2)
        fmt.Println("1 2")
    }
}
