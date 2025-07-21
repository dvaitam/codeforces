package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var A1, B1, C1 int
    var A2, B2, C2 int
    if _, err := fmt.Fscan(in, &A1, &B1, &C1); err != nil {
        return
    }
    if _, err := fmt.Fscan(in, &A2, &B2, &C2); err != nil {
        return
    }
    // Determine line types: 0 = empty, 1 = all points, 2 = normal line
    state := func(A, B, C int) int {
        if A == 0 && B == 0 {
            if C == 0 {
                return 1
            }
            return 0
        }
        return 2
    }
    s1 := state(A1, B1, C1)
    s2 := state(A2, B2, C2)
    // If either is empty, intersection is empty
    if s1 == 0 || s2 == 0 {
        fmt.Println(0)
        return
    }
    // If either is all points, intersection is the other set (if non-empty)
    if s1 != 2 || s2 != 2 {
        fmt.Println(-1)
        return
    }
    // Both are normal lines: check intersection
    det := A1*B2 - A2*B1
    if det != 0 {
        fmt.Println(1)
        return
    }
    // Parallel: check if same line
    if A1*C2 == A2*C1 && B1*C2 == B2*C1 {
        fmt.Println(-1)
    } else {
        fmt.Println(0)
    }
}
