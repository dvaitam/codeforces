package main

import "fmt"

func main() {
    var a, m int
    if _, err := fmt.Scan(&a, &m); err != nil {
        return
    }
    // Remove all factors of 2 from m
    g := m
    for g%2 == 0 {
        g /= 2
    }
    // Production stops iff the odd part of m divides a
    if a%g == 0 {
        fmt.Println("Yes")
    } else {
        fmt.Println("No")
    }
}
