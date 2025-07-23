package main

import (
    "fmt"
)

func main() {
    var n int64
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    var ans int64
    var digit int64 = 1
    var start int64 = 1
    // Sum digits for numbers with lengths less than that of n
    for start*10 <= n {
        count := (start*10 - start)
        ans += count * digit
        digit++
        start *= 10
    }
    // Remaining numbers with the same digit length as n
    ans += (n - start + 1) * digit
    fmt.Println(ans)
}
