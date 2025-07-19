package main

import "fmt"

func main() {
    var y, b, r int
    if _, err := fmt.Scan(&y, &b, &r); err != nil {
        return
    }
    m1 := 3*r - 3
    m2 := 3*b
    m3 := 3*y + 3
    res := m1
    if m2 < res {
        res = m2
    }
    if m3 < res {
        res = m3
    }
    fmt.Println(res)
}
