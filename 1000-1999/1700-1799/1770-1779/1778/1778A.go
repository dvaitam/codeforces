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

    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int, n)
        sum := 0
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
            sum += a[i]
        }
        hasNegPair := false
        hasOppPair := false
        for i := 0; i+1 < n; i++ {
            if a[i] == -1 && a[i+1] == -1 {
                hasNegPair = true
            }
            if a[i] != a[i+1] {
                hasOppPair = true
            }
        }
        if hasNegPair {
            fmt.Fprintln(out, sum+4)
        } else if hasOppPair {
            fmt.Fprintln(out, sum)
        } else {
            fmt.Fprintln(out, sum-4)
        }
    }
}

