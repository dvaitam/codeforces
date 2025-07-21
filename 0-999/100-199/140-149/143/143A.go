package main

import (
    "fmt"
)

func main() {
    var r1, r2, c1, c2, d1, d2 int
    // read row sums, column sums, diagonal sums
    if _, err := fmt.Scan(&r1, &r2, &c1, &c2, &d1, &d2); err != nil {
        return
    }
    // Try all possibilities for 2x2 square entries a, b, c, d
    // positions:
    // a b
    // c d
    for a := 1; a <= 9; a++ {
        for b := 1; b <= 9; b++ {
            if b == a || a+b != r1 {
                continue
            }
            for c := 1; c <= 9; c++ {
                if c == a || c == b {
                    continue
                }
                // check column1 sum
                if a+c != c1 {
                    continue
                }
                // compute d from second row sum
                d := r2 - c
                if d < 1 || d > 9 || d == a || d == b || d == c {
                    continue
                }
                // check column2 sum
                if b+d != c2 {
                    continue
                }
                // check main diagonal sum
                if a+d != d1 {
                    continue
                }
                // check anti-diagonal sum
                if b+c != d2 {
                    continue
                }
                // found valid assignment
                fmt.Printf("%d %d\n%d %d\n", a, b, c, d)
                return
            }
        }
    }
    // no solution
    fmt.Println(-1)
}
