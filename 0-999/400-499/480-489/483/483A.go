package main

import (
    "fmt"
)

func main() {
    var l, r int64
    if _, err := fmt.Scan(&l, &r); err != nil {
        return
    }
    // If fewer than 3 numbers, or exactly three and starting odd, no solution
    if r-l+1 < 3 || (r-l+1 == 3 && l%2 != 0) {
        fmt.Println(-1)
        return
    }
    // Choose first even start if possible
    if l%2 == 0 {
        fmt.Println(l, l+1, l+2)
    } else {
        // l is odd and segment length >3, so l+3 <= r
        fmt.Println(l+1, l+2, l+3)
    }
}
