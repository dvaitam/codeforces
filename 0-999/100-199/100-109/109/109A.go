package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    s := ""
    // Subtract 4 until remaining is divisible by 7
    for n%7 != 0 && n-4 >= 0 {
        n -= 4
        s += "4"
    }
    // If remaining is not divisible by 7, no solution
    if n%7 != 0 {
        fmt.Print(-1)
        return
    }
    // Print all the 4s
    fmt.Print(s)
    // Print the required number of 7s
    for i := 0; i < n/7; i++ {
        fmt.Print("7")
    }
}
