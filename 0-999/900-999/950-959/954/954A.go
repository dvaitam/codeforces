package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    var s string
    fmt.Fscan(reader, &s)
    res := 0
    for i := 0; i < n; {
        res++
        if i+1 < n && s[i] != s[i+1] {
            i += 2
        } else {
            i++
        }
    }
    fmt.Print(res)
}
