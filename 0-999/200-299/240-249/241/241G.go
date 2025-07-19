package main

import "fmt"

func main() {
    // Output count of balloons
    fmt.Println("302 ")
    // First balloon: position 0, max pressure 800000
    fmt.Println(" 0 800000")
    // Generate remaining balloons
    s := 60000
    for j := 300; j > 0; j-- {
        fmt.Printf("%d %d\n", s, j)
        s += 2*j - 1
    }
    // Last balloon
    fmt.Printf("%d %d\n", s+(1<<17), 800000)
}
