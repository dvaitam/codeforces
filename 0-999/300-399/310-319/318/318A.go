package main

import "fmt"

func main() {
    var n, k int64
    if _, err := fmt.Scan(&n, &k); err != nil {
        return
    }
    // Number of odd numbers from 1 to n
    half := (n + 1) / 2
    var ans int64
    if k <= half {
        ans = 2*k - 1
    } else {
        ans = 2 * (k - half)
    }
    fmt.Println(ans)
}
