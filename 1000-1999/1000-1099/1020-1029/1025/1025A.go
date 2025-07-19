package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    var s string
    if _, err := fmt.Scan(&s); err != nil {
        return
    }
    // Single character is always trivially repeated
    if n == 1 {
        fmt.Println("Yes")
        return
    }
    visited := make([]bool, 128)
    for i := 0; i < n && i < len(s); i++ {
        c := s[i]
        if visited[c] {
            fmt.Println("Yes")
            return
        }
        visited[c] = true
    }
    fmt.Println("No")
}
