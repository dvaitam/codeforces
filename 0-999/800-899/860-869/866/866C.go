package main

import "fmt"

func main() {
    var s string
    if _, err := fmt.Scan(&s); err != nil {
        return
    }
    n := len(s)
    for i := 0; i < n/2; i++ {
        if s[i] != s[n-1-i] {
            fmt.Println("NO")
            return
        }
    }
    fmt.Println("YES")
}
