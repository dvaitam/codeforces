package main

import (
    "bufio"
    "fmt"
    "os"
)

// This program outputs the k-th positive integer that is not divisible by 3
// and does not end with the digit 3. We precompute the first 1000 such
// integers and answer each query directly.
func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    // Precompute the sequence up to 1000 terms
    seq := make([]int, 0, 1000)
    for x := 1; len(seq) < 1000; x++ {
        if x%3 != 0 && x%10 != 3 {
            seq = append(seq, x)
        }
    }

    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var k int
        fmt.Fscan(in, &k)
        if k >= 1 && k <= len(seq) {
            fmt.Fprintln(out, seq[k-1])
        }
    }
}
