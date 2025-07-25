package main

import "fmt"

func main() {
    var n, d int
    if _, err := fmt.Scan(&n, &d); err != nil {
        return
    }
    var m int
    fmt.Scan(&m)
    for i := 0; i < m; i++ {
        var x, y int
        fmt.Scan(&x, &y)
        if x+y >= d && x+y <= 2*n-d && y-x >= -d && y-x <= d {
            fmt.Println("YES")
        } else {
            fmt.Println("NO")
        }
    }
}
