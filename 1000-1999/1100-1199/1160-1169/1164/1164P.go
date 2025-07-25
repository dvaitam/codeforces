package main

import "fmt"

// This program solves the following problem:
// One draws fifteen lines in the plane. What is the largest number of
// equilateral triangles (whose sides lie on these lines) that one can generate?
// The optimal arrangement is three families of five parallel lines at 60°
// to each other, yielding 5 * 5 * 5 = 125 equilateral triangles.
func main() {
    fmt.Println(125)
}
