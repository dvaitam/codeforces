package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    // Read input string with exactly 6 decimal places
    var s string
    if _, err := fmt.Fscan(os.Stdin, &s); err != nil {
        return
    }
    // Try all p/q where 1 <= p, q <= 10
    for p := 1; p <= 10; p++ {
        for q := 1; q <= 10; q++ {
            val := float64(p) / float64(q)
            // Format to 6 decimal places
            if strconv.FormatFloat(val, 'f', 6, 64) == s {
                fmt.Printf("%d %d", p, q)
                return
            }
        }
    }
}
