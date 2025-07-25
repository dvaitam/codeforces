package main

import "fmt"

// This program solves the following problem:
// Count the number of 7-digit permutations using digits 1..7 exactly once
// such that digits 1 and 2 are not adjacent.
// Total permutations = 7! = 5040
// Permutations with 1 and 2 adjacent: treat as a block (2!*6! = 1440)
// Answer = 5040 - 1440 = 3600
func main() {
    fmt.Println(5040 - 1440)
}
