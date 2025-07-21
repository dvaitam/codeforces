package main

import "fmt"

func main() {
    var s string
    if _, err := fmt.Scan(&s); err != nil {
        return
    }
    b := []byte(s)
    for i := range b {
        if b[i] >= 'a' && b[i] <= 'z' {
            b[i] = b[i] - 'a' + 'A'
        }
    }
    fmt.Println(string(b))
}
