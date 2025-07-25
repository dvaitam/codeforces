package main

import "fmt"

// Real numbers a and b satisfy (a + b - 1)^2 = ab + 1.
// Let u = a + b and v = ab. The constraint gives v = (u - 1)^2 - 1 = u^2 - 2u.
// Then a^2 + b^2 = u^2 - 2v = -u^2 + 4u, which attains its maximum when u = 2.
// Maximum value = -(2)^2 + 4*2 = 4.
func main() {
    fmt.Println(4)
}
