package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n, k int
    var s string
    fmt.Fscan(reader, &n, &k, &s)

    cnt, top := 0, 0
    half := k / 2
    for i := 0; i < n; i++ {
        if s[i] == '(' && cnt < half {
            writer.WriteByte('(')
            cnt++
            top++
        } else if s[i] == ')' && top > 0 {
            writer.WriteByte(')')
            top--
        }
    }
}
