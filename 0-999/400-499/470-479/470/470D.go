package main

import "fmt"

func main() {
    var k int
    var s string
    if _, err := fmt.Scan(&k); err != nil {
        return
    }
    if _, err := fmt.Scan(&s); err != nil {
        return
    }
    res := make([]byte, len(s))
    shift := byte(k % 26)
    for i := 0; i < len(s); i++ {
        ch := s[i] - 'A'
        ch = (ch + shift) % 26
        res[i] = 'A' + ch
    }
    fmt.Println(string(res))
}
