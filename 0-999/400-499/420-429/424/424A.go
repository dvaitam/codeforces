package main

import (
    "fmt"
)

func main() {
    var n int
    fmt.Scan(&n)
    var str string
    fmt.Scan(&str)
    b := []byte(str)
    desired := n / 2
    cur := 0
    for _, c := range b {
        if c == 'X' {
            cur++
        }
    }
    moves := desired - cur
    if moves < 0 {
        moves = -moves
    }
    // Adjust the string to have exactly desired 'X's
    need := desired - cur
    if need > 0 {
        for i := range b {
            if need == 0 {
                break
            }
            if b[i] == 'x' {
                b[i] = 'X'
                need--
            }
        }
    } else if need < 0 {
        for i := range b {
            if need == 0 {
                break
            }
            if b[i] == 'X' {
                b[i] = 'x'
                need++
            }
        }
    }
    fmt.Println(moves)
    fmt.Println(string(b))
}
