package main

import (
    "fmt"
)

func main() {
    var x, y, n int64
    const mod = 1000000007
    if _, err := fmt.Scan(&x, &y); err != nil {
        return
    }
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    // sequence: f1=x, f2=y, f3=y-x, f4=-x, f5=-y, f6=x-y, then repeats
    r := n % 6
    var res int64
    switch r {
    case 1:
        res = x
    case 2:
        res = y
    case 3:
        res = y - x
    case 4:
        res = -x
    case 5:
        res = -y
    case 0:
        res = x - y
    }
    res %= mod
    if res < 0 {
        res += mod
    }
    fmt.Println(res)
}
