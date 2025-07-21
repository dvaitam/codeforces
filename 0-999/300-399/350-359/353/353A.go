package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    sumUp, sumDown := 0, 0
    swapPossible := false
    for i := 0; i < n; i++ {
        var x, y int
        fmt.Scan(&x, &y)
        sumUp += x
        sumDown += y
        if (x%2) != (y%2) {
            swapPossible = true
        }
    }
    // If both sums are even, no swaps needed
    if sumUp%2 == 0 && sumDown%2 == 0 {
        fmt.Println(0)
    } else if sumUp%2 == 1 && sumDown%2 == 1 && swapPossible {
        // Both odd and there is a piece with different parity halves
        fmt.Println(1)
    } else {
        // Impossible otherwise
        fmt.Println(-1)
    }
}
