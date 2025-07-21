package main

import "fmt"

func main() {
    var a, b, c int
    if _, err := fmt.Scan(&a, &b, &c); err != nil {
        return
    }
    // a > b ensured by problem constraints
    delta := a - b
    // minimal waiting time t satisfies b*t >= delta*c
    numerator := delta * c
    t := (numerator + b - 1) / b
    fmt.Println(t)
}
