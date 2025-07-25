package main

import (
    "bufio"
    "fmt"
    "os"
)

func isPalindrome(s string) bool {
    i, j := 0, len(s)-1
    for i < j {
        if s[i] != s[j] {
            return false
        }
        i++
        j--
    }
    return true
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n, k int
        var s string
        fmt.Fscan(in, &n, &k)
        fmt.Fscan(in, &s)
        if k == 0 || isPalindrome(s) {
            fmt.Fprintln(out, 1)
        } else {
            fmt.Fprintln(out, 2)
        }
    }
}

