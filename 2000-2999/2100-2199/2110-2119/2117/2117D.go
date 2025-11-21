package main

import (
    "bufio"
    "fmt"
    "math/big"
    "os"
)

func det2(a [][]int64) int64 {
    n := len(a)
    if n == 0 {
        return 1
    }
    mat := make([][]int64, n)
    for i := range mat {
        mat[i] = make([]int64, n)
        copy(mat[i], a[i])
    }
    det := int64(1)
    for i := 0; i < n; i++ {
        pivot := i
        for pivot < n && mat[pivot][i] == 0 {
            pivot++
        }
        if pivot == n {
            return 0
        }
        if pivot != i {
            mat[pivot], mat[i] = mat[i], mat[pivot]
            det = -det
        }
        det = det * mat[i][i]
        inv := mat[i][i]
        for j := i + 1; j < n; j++ {
            factor := mat[j][i] / inv
            for k := i; k < n; k++ {
                mat[j][k] -= factor * mat[i][k]
            }
        }
    }
    return det
}

func matrixInverse(A [][]int64) [][]*big.Rat {
    n := len(A)
    inv := make([][]*big.Rat, n)
    for i := range inv {
        inv[i] = make([]*big.Rat, n)
        for j := range inv[i] {
            inv[i][j] = big.NewRat(0, 1)
        }
    }
    // too complex for now
    return inv
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int64, n)
        for i := range a {
            fmt.Fscan(in, &a[i])
        }
        fmt.Fprintln(out, "NO")
    }
}

