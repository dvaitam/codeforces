package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod = 1000000007

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    fmt.Fscan(reader, &t)
    for ; t > 0; t-- {
        var n int
        var l, r int
        fmt.Fscan(reader, &n, &l, &r)
        // TODO: implement the algorithm for counting excellent arrays.
        fmt.Fprintln(writer, 0)
    }
}
