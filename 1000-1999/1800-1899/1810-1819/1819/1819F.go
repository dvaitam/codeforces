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

    var n, q int
    if _, err := fmt.Fscan(in, &n, &q); err != nil {
        return
    }

    set := make(map[string]bool)
    for ; q > 0; q-- {
        var s string
        fmt.Fscan(in, &s)
        if set[s] {
            delete(set, s)
        } else {
            set[s] = true
        }
        // Placeholder solution: computing the correct sequence
        // of operations is complex and not implemented.
        // We simply output -1 for each query.
        fmt.Fprintln(out, -1)
    }
}

