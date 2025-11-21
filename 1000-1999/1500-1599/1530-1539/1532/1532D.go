package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }
    skills := make([]int, n)
    for i := range skills {
        fmt.Fscan(in, &skills[i])
    }
    sort.Ints(skills)
    total := 0
    for i := 0; i < n; i += 2 {
        total += skills[i+1] - skills[i]
    }
    fmt.Fprintln(out, total)
}
