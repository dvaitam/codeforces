package main

import "fmt"

// sumDigits returns the sum of digits of n.
func sumDigits(n int) int {
    sum := 0
    for n > 0 {
        sum += n % 10
        n /= 10
    }
    return sum
}

func main() {
    var a int
    if _, err := fmt.Scan(&a); err != nil {
        return
    }
    for n := a; ; n++ {
        if sumDigits(n)%4 == 0 {
            fmt.Println(n)
            break
        }
    }
}
