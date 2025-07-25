package main

import (
    "fmt"
)

func main() {
    // Compute minimal largest segment length for 11 segments with no three forming a triangle
    // This yields Fibonacci sequence: a1=1, a2=1, ak = a(k-1) + a(k-2), so a11 = 89
    fmt.Println(89)
}
