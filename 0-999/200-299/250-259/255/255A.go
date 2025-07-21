package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    sums := [3]int{}
    for i := 0; i < n; i++ {
        var x int
        if _, err := fmt.Scan(&x); err != nil {
            return
        }
        sums[i%3] += x
    }
    // Determine which muscle has the maximum total repetitions
    if sums[0] > sums[1] && sums[0] > sums[2] {
        fmt.Println("chest")
    } else if sums[1] > sums[0] && sums[1] > sums[2] {
        fmt.Println("biceps")
    } else {
        fmt.Println("back")
    }
}
