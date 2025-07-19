package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    freq := make(map[int]int)
    var x int
    maxCount := 0
    for i := 0; i < n; i++ {
        fmt.Scan(&x)
        freq[x]++
        if freq[x] > maxCount {
            maxCount = freq[x]
        }
    }
    fmt.Println(maxCount)
}
