package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    for i := 1; i <= n; i++ {
        for j := 0; j < i; j++ {
            fmt.Print("*")
        }
        fmt.Println()
    }
}
