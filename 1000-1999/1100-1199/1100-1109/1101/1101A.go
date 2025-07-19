package main

import "fmt"

func main() {
    var q int
    if _, err := fmt.Scan(&q); err != nil {
        return
    }
    for i := 0; i < q; i++ {
        var l, r, d int64
        fmt.Scan(&l, &r, &d)
        if d < l {
            fmt.Println(d)
        } else {
            fmt.Println((r/d + 1) * d)
        }
    }
}
