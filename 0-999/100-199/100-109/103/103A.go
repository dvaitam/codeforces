package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    a := make([]int64, n)
    for i := 0; i < n; i++ {
        fmt.Scan(&a[i])
    }
    var ans int64 = int64(n)
    for i, v := range a {
        ans += (v - 1) * int64(i+1)
    }
    fmt.Println(ans)
}
