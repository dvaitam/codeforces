package main

import "fmt"

func main() {
    var a, b, c int
    var d, e, f int
    if _, err := fmt.Scan(&a, &b, &c); err != nil {
        return
    }
    if _, err := fmt.Scan(&d, &e, &f); err != nil {
        return
    }
    if a > d || a+b > d+e || a+b+c > d+e+f {
        fmt.Println("NO")
    } else {
        fmt.Println("YES")
    }
}
