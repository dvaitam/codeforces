package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
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
        arr := make([]int, n)
        hasOdd, hasEven := false, false
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &arr[i])
            if arr[i]%2 == 0 {
                hasEven = true
            } else {
                hasOdd = true
            }
        }
        if hasOdd && hasEven {
            sort.Ints(arr)
        }
        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, arr[i])
        }
        fmt.Fprintln(out)
    }
}
