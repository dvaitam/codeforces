package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    m := n - 10
    switch {
    case m < 1 || m > 11:
        fmt.Println(0)
    case m == 10:
        fmt.Println(15)
    default:
        fmt.Println(4)
    }
}
