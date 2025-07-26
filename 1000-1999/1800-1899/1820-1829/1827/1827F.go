package main

import (
    "bufio"
    "fmt"
    "os"
)

// TODO: implement the algorithm for Codeforces problem 1827F "Copium Permutation".
// The official statement can be found in problemF.txt. This is a placeholder
// solution that reads the input format and outputs zeros so that the program
// compiles and can be used as a template for a future full implementation.

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n int
        if _, err := fmt.Fscan(reader, &n); err != nil {
            return
        }
        arr := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(reader, &arr[i])
        }

        // The correct algorithm is not yet implemented. We simply output
        // n+1 zeros, one for each value of k in [0, n].
        for i := 0; i <= n; i++ {
            if i > 0 {
                writer.WriteByte(' ')
            }
            fmt.Fprint(writer, 0)
        }
        writer.WriteByte('\n')
    }
}
