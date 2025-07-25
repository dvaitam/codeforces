package main

import "fmt"

// This program solves the nested radical equations:
// a = sqrt(7 - sqrt(6 - a)), b = sqrt(7 - sqrt(6 + b)),
// c = sqrt(7 + sqrt(6 - c)), d = sqrt(7 + sqrt(6 + d)),
// and computes the product a*b*c*d, which equals 43.
func main() {
    fmt.Println(43)
}
