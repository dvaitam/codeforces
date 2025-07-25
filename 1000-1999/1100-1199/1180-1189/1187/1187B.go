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

    var n int
    fmt.Fscan(reader, &n)
    var s string
    fmt.Fscan(reader, &s)

    // Precompute positions of each letter in s
    pos := make([][]int, 26)
    for i := 0; i < n; i++ {
        c := s[i] - 'a'
        pos[c] = append(pos[c], i+1)
    }

    var m int
    fmt.Fscan(reader, &m)
    for i := 0; i < m; i++ {
        var t string
        fmt.Fscan(reader, &t)
        // Count occurrences and track max position needed
        cnt := [26]int{}
        ans := 0
        for j := 0; j < len(t); j++ {
            c := t[j] - 'a'
            cnt[c]++
            // position of the cnt[c]-th occurrence of c
            p := pos[c][cnt[c]-1]
            if p > ans {
                ans = p
            }
        }
        fmt.Fprintln(writer, ans)
    }
}
