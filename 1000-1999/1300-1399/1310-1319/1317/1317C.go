package main

import "fmt"

const mod = int64(1000000007)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    a, b := int64(0), int64(1)
    for i := 0; i < n; i++ {
        a, b = b, (a+b)%mod
    }
    fmt.Println(a)
}
