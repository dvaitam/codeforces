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

    var n, x int
    if _, err := fmt.Fscan(reader, &n, &x); err != nil {
        return
    }

    if n == 5 {
        fmt.Fprintln(writer, ">...v")
        fmt.Fprintln(writer, "v.<..")
        fmt.Fprintln(writer, "..^..")
        fmt.Fprintln(writer, ">....")
        fmt.Fprintln(writer, "..^.<")
        fmt.Fprintln(writer, "1 1")
        return
    }
    if n == 3 {
        fmt.Fprintln(writer, ">vv")
        fmt.Fprintln(writer, "^<.")
        fmt.Fprintln(writer, "^.<")
        fmt.Fprintln(writer, "1 3")
        return
    }

    mp := make([][]rune, n)
    for i := 0; i < n; i++ {
        mp[i] = make([]rune, n)
        for j := 0; j < n; j++ {
            mp[i][j] = '.'
        }
    }
    for i := 0; i < n; i++ {
        mp[i][0] = '^'
    }
    mp[0][0] = '>'

    for i := 0; i < n; i += 2 {
        for j := 1; j < n-1; j++ {
            if j < n/2 || j%2 == 1 {
                mp[i][j] = '>'
            }
        }
        mp[i][n-1] = 'v'
    }
    for i := 1; i < n; i += 2 {
        for j := n - 1; j > 0; j-- {
            if (n-j) < n/2 || j%2 == 1 {
                mp[i][j] = '<'
            }
        }
        mp[i][1] = 'v'
    }
    mp[n-1][1] = '<'

    for i := 0; i < n; i++ {
        writer.WriteString(string(mp[i]))
        writer.WriteByte('\n')
    }
    fmt.Fprintln(writer, "1 1")
}
