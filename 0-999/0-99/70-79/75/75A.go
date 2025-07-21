package main

import (
    "fmt"
)

// removeZeros removes all zero digits from n and returns the resulting integer.
func removeZeros(n int) int {
    res := 0
    place := 1
    for n > 0 {
        d := n % 10
        if d != 0 {
            res = d*place + res
            place *= 10
        }
        n /= 10
    }
    return res
}

func main() {
    var a, b int
    // Read two integers a and b
    fmt.Scan(&a)
    fmt.Scan(&b)
    c := a + b
    // Check if after removing zeros the equation remains valid
    if removeZeros(a)+removeZeros(b) == removeZeros(c) {
        fmt.Println("YES")
    } else {
        fmt.Println("NO")
    }
}
