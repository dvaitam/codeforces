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

    var b int
    if _, err := fmt.Fscan(reader, &b); err != nil {
        return
    }
    for i := 1; i <= b; i++ {
        for k := 0; k < b; k++ {
            var val int
            if k < b/2 {
                val = k*b + i
            } else {
                val = k*b + b - (i - 1)
            }
            fmt.Fprintf(writer, "%d ", val)
        }
        fmt.Fprint(writer, "\n")
    }
}
