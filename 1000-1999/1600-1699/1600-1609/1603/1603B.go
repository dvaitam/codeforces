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
    for i := 0; i < t; i++ {
        var x, y int64
        fmt.Fscan(in, &x, &y)
        if x > y {
            fmt.Fprintln(out, x+y)
        } else {
            x2 := x + ((y-x)/x)*x
            fmt.Fprintln(out, (x2+y)/2)
        }
    }
}
