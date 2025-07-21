package main

import (
    "fmt"
)

func main() {
    var n, m int
    if _, err := fmt.Scan(&n, &m); err != nil {
        return
    }
    rem := m
    i := 1
    for {
        if rem < i {
            break
        }
        rem -= i
        i++
        if i > n {
            i = 1
        }
    }
    fmt.Println(rem)
}
