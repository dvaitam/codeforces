package main

import (
    "fmt"
)

func main() {
    var a, b, c, d, e, f int64
    if _, err := fmt.Scan(&a, &b, &c, &d, &e, &f); err != nil {
        return
    }
    // Infinite gold from spell2 directly
    if c == 0 && d > 0 {
        fmt.Println("Ron")
        return
    }
    // Infinite lead -> gold
    if a == 0 && b > 0 && ((c > 0 && d > 0) || (c == 0 && d > 0)) {
        fmt.Println("Ron")
        return
    }
    // Infinite sand -> lead -> gold
    if e == 0 && f > 0 && ( (a > 0 && b > 0) && ((c > 0 && d > 0) || (c == 0 && d > 0)) ) {
        fmt.Println("Ron")
        return
    }
    // Positive cycle S->L->G->S
    if a > 0 && b > 0 && c > 0 && d > 0 && e > 0 && f > 0 {
        // Check if b/a * d/c * f/e > 1 using int arithmetic
        if b*d*f > a*c*e {
            fmt.Println("Ron")
            return
        }
    }
    fmt.Println("Hermione")
}
