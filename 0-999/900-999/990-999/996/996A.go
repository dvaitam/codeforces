package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    bills := []int{100, 20, 10, 5, 1}
    count := 0
    for _, b := range bills {
        count += n / b
        n %= b
    }
    fmt.Println(count)
}
