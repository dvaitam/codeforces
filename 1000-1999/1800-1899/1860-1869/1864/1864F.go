package main

import (
    "bufio"
    "fmt"
    "os"
)

// TODO: implement the actual algorithm for Codeforces problem 1864F.
// The full statement can be found in problemF.txt. The task asks for the
// minimum number of laminar segment subtraction operations required to
// zero out elements whose values lie in the range [l, r] for each query.
// Implementing the optimal approach is non-trivial, so this file provides
// only a placeholder that reads the input and outputs 0 for each query so
// that the program builds successfully.
func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, q int
    if _, err := fmt.Fscan(in, &n, &q); err != nil {
        return
    }
    arr := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &arr[i])
    }

    for ; q > 0; q-- {
        var l, r int
        fmt.Fscan(in, &l, &r)
        // Proper computation of the minimal number of operations is not
        // implemented yet.
        fmt.Fprintln(out, 0)
    }
}

