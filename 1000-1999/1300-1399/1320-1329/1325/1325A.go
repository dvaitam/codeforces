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

    var t, x int
    fmt.Fscan(reader, &t)
    for i := 0; i < t; i++ {
        fmt.Fscan(reader, &x)
        fmt.Fprintln(writer, 1, x-1)
    }
}
