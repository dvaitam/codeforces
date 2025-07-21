package main

import "fmt"

func main() {
    var n int
    var s string
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    if _, err := fmt.Scan(&s); err != nil {
        return
    }
    candies := make([]int, n)
    for i := range candies {
        candies[i] = 1
    }
    // forward pass
    for i := 1; i < n; i++ {
        switch s[i-1] {
        case 'R':
            candies[i] = candies[i-1] + 1
        case '=':
            candies[i] = candies[i-1]
        }
    }
    // backward pass
    for i := n - 2; i >= 0; i-- {
        switch s[i] {
        case 'L':
            if candies[i] <= candies[i+1] {
                candies[i] = candies[i+1] + 1
            }
        case '=':
            if candies[i] != candies[i+1] {
                candies[i] = candies[i+1]
            }
        }
    }
    // output result
    for i, c := range candies {
        if i > 0 {
            fmt.Print(" ")
        }
        fmt.Print(c)
    }
    fmt.Println()
}
