package main

import "fmt"

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
    if b == 0 {
        return a
    }
    return gcd(b, a%b)
}

func main() {
    var Y, W int
    if _, err := fmt.Scan(&Y, &W); err != nil {
        return
    }
    // Dot wins if she rolls at least max(Y, W)
    m := Y
    if W > m {
        m = W
    }
    // favorable outcomes: rolls from m to 6 inclusive
    num := 7 - m
    den := 6
    // simplify fraction
    g := gcd(num, den)
    num /= g
    den /= g
    fmt.Printf("%d/%d", num, den)
}
