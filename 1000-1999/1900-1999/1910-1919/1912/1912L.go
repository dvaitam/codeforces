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

    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }
    var s string
    fmt.Fscan(in, &s)

    prefL := make([]int, n+1)
    for i := 0; i < n; i++ {
        prefL[i+1] = prefL[i]
        if s[i] == 'L' {
            prefL[i+1]++
        }
    }
    totalL := prefL[n]
    totalO := n - totalL

    for k := 1; k < n; k++ {
        leftL := prefL[k]
        leftO := k - leftL
        rightL := totalL - leftL
        rightO := totalO - leftO
        if leftL != rightL && leftO != rightO {
            fmt.Fprintln(out, k)
            return
        }
    }
    fmt.Fprintln(out, -1)
}
