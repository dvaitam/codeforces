package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    count := 0
    for i := 0; i < n; i++ {
        var p, q int
        if _, err := fmt.Scan(&p, &q); err != nil {
            return
        }
        if q-p >= 2 {
            count++
        }
    }
    fmt.Println(count)
}
