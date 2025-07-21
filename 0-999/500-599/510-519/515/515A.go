package main

import "fmt"

func main() {
    var a, b, s int64
    _, err := fmt.Scan(&a, &b, &s)
    if err != nil {
        return
    }
    d := abs(a) + abs(b)
    if s >= d && (s-d)%2 == 0 {
        fmt.Println("Yes")
    } else {
        fmt.Println("No")
    }
}

func abs(x int64) int64 {
    if x < 0 {
        return -x
    }
    return x
}
