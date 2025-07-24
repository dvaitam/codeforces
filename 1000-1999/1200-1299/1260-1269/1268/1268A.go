package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, k int
    if _, err := fmt.Fscan(in, &n, &k); err != nil {
        return
    }
    var s string
    fmt.Fscan(in, &s)

    b := []byte(s)
    for i := k; i < n; i++ {
        b[i] = b[i-k]
    }

    if ge(b, []byte(s)) {
        fmt.Fprintln(out, n)
        fmt.Fprintln(out, string(b))
        return
    }

    // increment the first k digits
    carry := byte(1)
    for i := k - 1; i >= 0 && carry > 0; i-- {
        d := b[i] - '0' + carry
        b[i] = d%10 + '0'
        carry = d / 10
    }

    for i := k; i < n; i++ {
        b[i] = b[i-k]
    }

    fmt.Fprintln(out, n)
    fmt.Fprintln(out, string(b))
}

func ge(a, b []byte) bool {
    for i := range a {
        if a[i] > b[i] {
            return true
        }
        if a[i] < b[i] {
            return false
        }
    }
    return true
}
