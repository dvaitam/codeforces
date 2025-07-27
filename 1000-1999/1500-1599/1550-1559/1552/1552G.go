package main

import (
    "bufio"
    "fmt"
    "os"
)

// TODO: implement a correct solution for problem G
func main() {
    in := bufio.NewReader(os.Stdin)
    var n, k int
    fmt.Fscan(in, &n, &k)
    for i := 0; i < k; i++ {
        var q int
        fmt.Fscan(in, &q)
        for j := 0; j < q; j++ {
            var x int
            fmt.Fscan(in, &x)
            _ = x
        }
    }
    fmt.Println("REJECTED")
}

