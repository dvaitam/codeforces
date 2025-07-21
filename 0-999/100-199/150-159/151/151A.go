package main

import (
    "fmt"
)

// main reads input parameters and computes the maximum number of toasts per friend.
func main() {
    var n, k, l, c, d, p, nl, np int
    // n: number of friends
    // k: number of bottles, each with l milliliters
    // c: number of limes, each cut into d slices
    // p: total grams of salt
    // nl: milliliters needed for one toast per person
    // np: grams of salt needed for one toast per person
    if _, err := fmt.Scan(&n, &k, &l, &c, &d, &p, &nl, &np); err != nil {
        return
    }
    totalDrink := k * l
    totalLimeSlices := c * d
    // maximum toasts limited by each resource
    toastsByDrink := totalDrink / nl
    toastsByLime := totalLimeSlices
    toastsBySalt := p / np
    // overall maximum toasts all friends can make
    maxToasts := min(toastsByDrink, min(toastsByLime, toastsBySalt))
    // toasts per friend
    fmt.Println(maxToasts / n)
}

// min returns the smaller of two integers.
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
