package main

import (
    "bufio"
    "fmt"
    "os"
)

func min(a int, b ...int) int {
    m := a
    for _, v := range b {
        if v < m {
            m = v
        }
    }
    return m
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var w, d, h int
        fmt.Fscan(reader, &w, &d, &h)
        var a, b, f, g int
        fmt.Fscan(reader, &a, &b, &f, &g)

        // consider four possible walls to move vertically along
        path1 := a + f + abs(b-g)
        path2 := (w - a) + (w - f) + abs(b-g)
        path3 := b + g + abs(a-f)
        path4 := (d - b) + (d - g) + abs(a-f)

        res := h + min(path1, path2, path3, path4)
        fmt.Fprintln(writer, res)
    }
}
