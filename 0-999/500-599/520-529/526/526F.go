package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    fmt.Fscan(in, &n)
    rows := make([]int, n)
    for i := 0; i < n; i++ {
        var r, c int
        fmt.Fscan(in, &r, &c)
        rows[r-1] = c
    }

    // Naive O(n^2) solution: iterate over all subarrays and
    // check if the set of columns forms a consecutive range.
    // This is too slow for the official constraints but serves
    // as a placeholder implementation.
    ans := 0
    for i := 0; i < n; i++ {
        minVal := rows[i]
        maxVal := rows[i]
        for j := i; j < n; j++ {
            if rows[j] < minVal {
                minVal = rows[j]
            }
            if rows[j] > maxVal {
                maxVal = rows[j]
            }
            if maxVal-minVal == j-i {
                ans++
            }
        }
    }
    fmt.Println(ans)
}

