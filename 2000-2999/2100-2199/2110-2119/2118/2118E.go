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

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n, m int
        fmt.Fscan(in, &n, &m)
        top, bottom := 0, n-1
        left, right := 0, m-1
        res := make([][2]int, 0, n*m)
        for top <= bottom && left <= right {
            for y := left; y <= right; y++ {
                res = append(res, [2]int{top, y})
            }
            top++
            if top > bottom {
                break
            }
            for x := top; x <= bottom; x++ {
                res = append(res, [2]int{x, right})
            }
            right--
            if left > right {
                break
            }
            for y := right; y >= left; y-- {
                res = append(res, [2]int{bottom, y})
            }
            bottom--
            if top > bottom {
                break
            }
            for x := bottom; x >= top; x-- {
                res = append(res, [2]int{x, left})
            }
            left++
        }
        for _, cell := range res {
            fmt.Fprintf(out, "%d %d\n", cell[0]+1, cell[1]+1)
        }
        fmt.Fprintln(out)
    }
}
