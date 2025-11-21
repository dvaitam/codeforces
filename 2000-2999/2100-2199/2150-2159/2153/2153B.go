package main

import (
    "bufio"
    "fmt"
    "os"
)

func possible(x, y, z int) bool {
    for x > 0 || y > 0 || z > 0 {
        ones := (x & 1) + (y & 1) + (z & 1)
        if ones == 2 {
            return false
        }
        x >>= 1
        y >>= 1
        z >>= 1
    }
    return true
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    fmt.Fscan(reader, &t)
    for ; t > 0; t-- {
        var x, y, z int
        fmt.Fscan(reader, &x, &y, &z)
        if possible(x, y, z) {
            fmt.Fprintln(writer, "YES")
        } else {
            fmt.Fprintln(writer, "NO")
        }
    }
}
