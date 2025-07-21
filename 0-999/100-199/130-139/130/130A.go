package main

import "fmt"

func main() {
    var n int64
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    // Hexagonal number formula: H_n = 2n^2 - n
    result := 2*n*n - n
    fmt.Println(result)
}
