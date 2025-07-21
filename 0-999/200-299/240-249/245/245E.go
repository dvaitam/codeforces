package main

import (
    "fmt"
)

func main() {
    var s string
    if _, err := fmt.Scan(&s); err != nil {
        return
    }
    cur := 0
    maxPeople := 0
    for _, c := range s {
        if c == '+' {
            cur++
            if cur > maxPeople {
                maxPeople = cur
            }
        } else if c == '-' {
            if cur > 0 {
                cur--
            } else {
                // someone leaves who was inside before start
                maxPeople++
            }
        }
    }
    fmt.Println(maxPeople)
}
