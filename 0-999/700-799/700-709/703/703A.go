package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    mishka, chris := 0, 0
    for i := 0; i < n; i++ {
        var m, c int
        fmt.Scan(&m, &c)
        if m > c {
            mishka++
        } else if c > m {
            chris++
        }
    }
    if mishka > chris {
        fmt.Println("Mishka")
    } else if chris > mishka {
        fmt.Println("Chris")
    } else {
        fmt.Println("Friendship is magic!^^")
    }
}
