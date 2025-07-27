package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    k := n / 2
    // Compute binomial coefficient C(n, k)
    var comb uint64 = 1
    for i := 1; i <= k; i++ {
        comb = comb * uint64(n-k+i) / uint64(i)
    }
    // Compute factorial of (k-1)
    var fact uint64 = 1
    for i := 1; i <= k-1; i++ {
        fact *= uint64(i)
    }
    // Number of ways: C(n, k)/2 * (k-1)! * (k-1)!
    result := comb * fact * fact / 2
    fmt.Println(result)
}
