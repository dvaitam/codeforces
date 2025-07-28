package main

import "fmt"

const mod = 1000000007

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    if n == 0 {
        fmt.Println(0)
        return
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, (a+b)%mod
    }
    fmt.Println(b)
}
