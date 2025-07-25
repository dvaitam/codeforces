package main

import "fmt"

func main() {
    var q int
    if _, err := fmt.Scan(&q); err != nil {
        return
    }
    for i := 0; i < q; i++ {
        var n, k int
        fmt.Scan(&n, &k)
        low, high := 1, 1000000000
        for j := 0; j < n; j++ {
            var a int
            fmt.Scan(&a)
            lo := a - k
            if lo > low {
                low = lo
            }
            hi := a + k
            if hi < high {
                high = hi
            }
        }
        if low > high {
            fmt.Println(-1)
        } else {
            fmt.Println(high)
        }
    }
}
