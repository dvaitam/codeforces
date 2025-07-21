package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    count := 0
    var prev string
    for i := 0; i < n; i++ {
        var s string
        fmt.Scan(&s)
        if i == 0 || s != prev {
            count++
        }
        prev = s
    }
    fmt.Println(count)
}
