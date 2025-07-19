package main

import (
    "fmt"
)

func main() {
    var n int64
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    tmp := n
    var a int64
    var ans int64
    for tmp >= 10 {
        tmp /= 10
        a = a*10 + 9
        ans += 9
    }
    b := n - a
    for b > 0 {
        ans += b % 10
        b /= 10
    }
    fmt.Println(ans)
}
