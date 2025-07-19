package main

import "fmt"

func main() {
    var n1, n2, n3, n4 int
    if _, err := fmt.Scan(&n1, &n2, &n3, &n4); err != nil {
        return
    }
    _ = n2
    var res int
    if n1 != n4 {
        res = 0
    } else {
        if n1 == 0 {
            if n3 != 0 {
                res = 0
            } else {
                res = 1
            }
        } else {
            res = 1
        }
    }
    fmt.Println(res)
}
