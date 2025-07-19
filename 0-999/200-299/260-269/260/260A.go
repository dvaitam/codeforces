package main

import (
    "fmt"
)

func main() {
    var a, b int64
    var n int
    if _, err := fmt.Scan(&a, &b, &n); err != nil {
        return
    }
    found := false
    var digit int
    for i := 0; i < 10; i++ {
        if (a*10+int64(i))%b == 0 {
            found = true
            digit = i
            break
        }
    }
    if !found {
        fmt.Println(-1)
        return
    }
    // output the number: a, the found digit, then n-1 zeros
    fmt.Print(a)
    fmt.Print(digit)
    for i := 0; i < n-1; i++ {
        fmt.Print(0)
    }
}
