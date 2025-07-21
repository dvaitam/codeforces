package main

import (
    "fmt"
)

func isPrime(x int) bool {
    if x < 2 {
        return false
    }
    for i := 2; i*i <= x; i++ {
        if x%i == 0 {
            return false
        }
    }
    return true
}

func main() {
    var n, m int
    _, err := fmt.Scan(&n, &m)
    if err != nil {
        return
    }
    for p := n + 1; ; p++ {
        if isPrime(p) {
            if p == m {
                fmt.Println("YES")
            } else {
                fmt.Println("NO")
            }
            break
        }
    }
}
