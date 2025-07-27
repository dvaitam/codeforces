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
        sum1 := 1 << n
        for i := 1; i <= n/2-1; i++ {
            sum1 += 1 << i
        }
        sum2 := 0
        for i := n / 2; i <= n-1; i++ {
            sum2 += 1 << i
        }
        fmt.Fprintln(out, sum1-sum2)
    }
}
