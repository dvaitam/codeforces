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

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var a, b, n int64
        fmt.Fscan(reader, &a, &b, &n)
        var res int64
        switch n % 3 {
        case 0:
            res = a
        case 1:
            res = b
        default:
            res = a ^ b
        }
        fmt.Fprintln(writer, res)
    }
}
