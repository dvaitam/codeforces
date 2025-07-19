package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    if n == 2 {
        fmt.Println(-1)
        return
    }
    fmt.Println(10)
    fmt.Println(15)
    for i := 2; i < n; i++ {
        fmt.Println(i*6 - 6)
    }
}
