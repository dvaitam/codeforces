package main

import (
    "fmt"
)

func main() {
    var n, m int
    if _, err := fmt.Scan(&n, &m); err != nil {
        return
    }
    // If only one card in total, probability is 1
    if n*m == 1 {
        fmt.Println("1")
        return
    }
    nf := float64(n)
    mf := float64(m)
    // Calculate probability: 1/n + ((n-1)/n)*((m-1)/(n*m-1))
    ans := 1.0/nf + ((nf-1.0)/nf)*((mf-1.0)/(nf*mf - 1.0))
    fmt.Printf("%.9f\n", ans)
}
