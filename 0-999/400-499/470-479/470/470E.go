package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if (i+j)%2 == 0 {
                fmt.Print(".")
            } else {
                fmt.Print("#")
            }
        }
        fmt.Println()
    }
}
