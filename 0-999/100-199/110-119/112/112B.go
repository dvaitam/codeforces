package main

import (
    "fmt"
)

func main() {
    var n, x, y int
    if _, err := fmt.Scan(&n, &x, &y); err != nil {
        return
    }
    // Special case: n == 2 is always NO
    if n == 2 {
        fmt.Print("NO")
        return
    }
    mid := n / 2
    // Check forbidden central positions
    if (x == mid && y == mid) ||
       (x == mid && y == mid+1) ||
       (x == mid+1 && y == mid) ||
       (x == mid+1 && y == mid+1) {
        fmt.Print("NO")
        return
    }
    fmt.Print("YES")
}
